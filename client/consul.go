/*
Copyright (C) 2020 Expedia, Inc.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package client

import (
	"fmt"

	"github.com/HotelsDotCom/go-logger"
	consul "github.com/hashicorp/consul/api"
)

type txnClient interface {
	Txn(consul.TxnOps, *consul.QueryOptions) (bool, *consul.TxnResponse, *consul.QueryMeta, error)
}

//Consul represents the consul client.
type Consul interface {
	KVTransact(datacenter string, operations []KVOperation) ([]KVTransactionResult, []KVTransactionError, error)
	IsVerbSupported(verb string) bool
}

type consulClient struct {
	txnClient txnClient
}

//NewConsul produces a new consul client
func NewConsul() (Consul, error) {
	client, err := consul.NewClient(consul.DefaultConfig())
	if nil != err {
		logger.Error("failed to initialize consul: %v", err)
		return nil, err
	}

	consul := &consulClient{
		txnClient: client.Txn(),
	}

	logger.Info("initialized consul")
	return consul, nil
}

func (c *consulClient) KVTransact(datacenter string, operations []KVOperation) ([]KVTransactionResult, []KVTransactionError, error) {
	var q *consul.QueryOptions = nil
	if "" != datacenter {
		q = &consul.QueryOptions{
			Datacenter: datacenter,
		}
	}
	input := consul.TxnOps{}
	for _, op := range operations {
		kvOp := toTxnKVOp(op)
		input = append(input, &consul.TxnOp{KV: &kvOp})
	}

	ok, response, _, err := c.txnClient.Txn(input, q)
	if nil != err {
		return nil, nil, fmt.Errorf("failed to make transaction request: %v", err)
	}
	if !ok {
		return nil, mapTxnErrorsToKVTransactionErrors(response.Errors, toKVTransactionError), nil
	}

	return mapTxnResultsToKVTransactionResults(response.Results), nil, nil
}

func (c *consulClient) IsVerbSupported(verb string) bool {
	_, retval := getSupportedVerbs()[verb]
	return retval
}

func mapTxnErrorsToKVTransactionErrors(source consul.TxnErrors, f func(*consul.TxnError) KVTransactionError) []KVTransactionError {
	target := make([]KVTransactionError, len(source))
	for index, value := range source {
		target[index] = f(value)
	}
	return target
}

func mapTxnResultsToKVTransactionResults(source consul.TxnResults) []KVTransactionResult {
	target := make([]KVTransactionResult, len(source))
	for index, value := range source {
		target[index] = KVTransactionResult{
			Index: index,
			Key:   value.KV.Key,
			Value: value.KV.Value,
		}
	}
	return target
}
