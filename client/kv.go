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
	"encoding/json"

	consul "github.com/hashicorp/consul/api"
)

//KVOperation represents a consul key-value operation.
type KVOperation struct {
	Verb  string          `json:"verb"`
	Key   string          `json:"key"`
	Value json.RawMessage `json:"value"`
}

//KVTransactionError represents a key-value transaction error.
type KVTransactionError struct {
	Index int    `json:"index"`
	Error string `json:"error"`
}

//KVTransactionResult represents a key-value transaction result.
type KVTransactionResult struct {
	Index int             `json:"index"`
	Key   string          `json:"key"`
	Value json.RawMessage `json:"value"`
}

func toTxnKVOp(kvOp KVOperation) consul.KVTxnOp {
	return consul.KVTxnOp{
		Verb:  getSupportedVerbs()[kvOp.Verb],
		Key:   kvOp.Key,
		Value: kvOp.Value,
	}
}

func toKVTransactionError(txnError *consul.TxnError) KVTransactionError {
	return KVTransactionError{
		Index: txnError.OpIndex,
		Error: txnError.What,
	}
}

func getSupportedVerbs() map[string]consul.KVOp {
	retval := make(map[string]consul.KVOp)
	retval[string(consul.KVSet)] = consul.KVSet
	retval[string(consul.KVDelete)] = consul.KVDelete
	retval[string(consul.KVDeleteCAS)] = consul.KVDeleteCAS
	retval[string(consul.KVDeleteTree)] = consul.KVDeleteTree
	retval[string(consul.KVCAS)] = consul.KVCAS
	retval[string(consul.KVLock)] = consul.KVLock
	retval[string(consul.KVUnlock)] = consul.KVUnlock
	retval[string(consul.KVGet)] = consul.KVGet
	retval[string(consul.KVGetTree)] = consul.KVGetTree
	retval[string(consul.KVCheckSession)] = consul.KVCheckIndex
	retval[string(consul.KVCheckIndex)] = consul.KVCheckIndex
	retval[string(consul.KVCheckNotExists)] = consul.KVCheckNotExists
	return retval
}
