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
	"errors"
	"testing"

	"github.com/HotelsDotCom/go-logger/loggertest"
	consul "github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var ConsulImpl Consul
var ConsulMockClient *MockClient

func Before(t *testing.T) {
	loggertest.Init("DEBUG")
	ConsulImpl, _ = NewConsul()
	ConsulMockClient = NewMockClient(t)
	ConsulImpl.(*consulClient).txnClient = ConsulMockClient
}

func After() {
	loggertest.Reset()
}

func TestKVTransactWithDC(t *testing.T) {
	Before(t)
	defer After()

	const key = "joe/mama"
	op := getOperation(key)
	ops := []KVOperation{op}

	ConsulMockClient.TxnFunc = func(operations consul.TxnOps, queryOptions *consul.QueryOptions) (bool, *consul.TxnResponse, *consul.QueryMeta, error) {
		require.NotNil(t, queryOptions)
		assert.Equal(t, "dc", queryOptions.Datacenter)
		response := getSuccessfulResponse(key)
		return true, &response, nil, nil
	}

	result, rollback, err := ConsulImpl.KVTransact("dc", ops)
	assert.Nil(t, err)
	assert.Nil(t, rollback)
	require.Equal(t, 1, len(result))

	assert.Equal(t, 0, result[0].Index)
	assert.Equal(t, key, result[0].Key)
	assert.Nil(t, result[0].Value)
}

func TestKVTransact(t *testing.T) {
	Before(t)
	defer After()

	const key = "joe/mama"
	op := getOperation(key)
	ops := []KVOperation{op}

	ConsulMockClient.TxnFunc = func(operations consul.TxnOps, queryOptions *consul.QueryOptions) (bool, *consul.TxnResponse, *consul.QueryMeta, error) {
		assert.Nil(t, queryOptions)
		response := getSuccessfulResponse(key)
		return true, &response, nil, nil
	}

	result, rollback, err := ConsulImpl.KVTransact("", ops)
	assert.Nil(t, err)
	assert.Nil(t, rollback)
	require.Equal(t, 1, len(result))

	assert.Equal(t, 0, result[0].Index)
	assert.Equal(t, key, result[0].Key)
	assert.Nil(t, result[0].Value)
}

func TestKVTransactFailed(t *testing.T) {
	Before(t)
	defer After()

	const key = "joe/mama"
	op := getOperation(key)
	ops := []KVOperation{op}

	ConsulMockClient.TxnFunc = func(operations consul.TxnOps, queryOptions *consul.QueryOptions) (bool, *consul.TxnResponse, *consul.QueryMeta, error) {
		return false, nil, nil, errors.New("kablammo")
	}

	result, rollback, err := ConsulImpl.KVTransact("", ops)
	assert.Nil(t, result)
	assert.Nil(t, rollback)
	require.NotNil(t, err)
	assert.Equal(t, "failed to make transaction request: kablammo", err.Error())
}

func TestKVTransactRollback(t *testing.T) {
	Before(t)
	defer After()

	const key = "joe/mama"
	const errorMessage = "kablammo"
	op := getOperation(key)
	ops := []KVOperation{op}

	ConsulMockClient.TxnFunc = func(operations consul.TxnOps, queryOptions *consul.QueryOptions) (bool, *consul.TxnResponse, *consul.QueryMeta, error) {
		response := getErrorResponse(errorMessage)
		return false, &response, nil, nil
	}

	result, rollback, err := ConsulImpl.KVTransact("", ops)
	assert.Nil(t, result)
	assert.Nil(t, err)
	require.Equal(t, 1, len(rollback))
	assert.Equal(t, 0, rollback[0].Index)
	assert.Equal(t, errorMessage, rollback[0].Error)
}

func TestIsVerbSupported(t *testing.T) {
	Before(t)
	defer After()

	assert.True(t, ConsulImpl.IsVerbSupported(string(consul.KVSet)))
	assert.False(t, ConsulImpl.IsVerbSupported("jump"))
}

func getOperation(key string) KVOperation {
	return KVOperation{
		Verb:  string(consul.KVSet),
		Key:   key,
		Value: []byte("{\"one\":1,\"two\":\"two\",\"three\":true,\"four\":[],\"five\":null}"),
	}
}

func getSuccessfulResponse(key string) consul.TxnResponse {
	return consul.TxnResponse{
		Results: consul.TxnResults{
			&consul.TxnResult{
				KV: &consul.KVPair{
					Key:   key,
					Value: nil,
				},
			},
		},
		Errors: consul.TxnErrors{},
	}
}

func getErrorResponse(errorMessage string) consul.TxnResponse {
	return consul.TxnResponse{
		Results: consul.TxnResults{},
		Errors: consul.TxnErrors{
			&consul.TxnError{
				OpIndex: 0,
				What:    errorMessage,
			},
		},
	}
}

type MockClient struct {
	t       *testing.T
	TxnFunc func(operations consul.TxnOps, queryOptions *consul.QueryOptions) (bool, *consul.TxnResponse, *consul.QueryMeta, error)
}

func NewMockClient(t *testing.T) *MockClient {
	m := &MockClient{t: t}
	m.TxnFunc = func(operations consul.TxnOps, queryOptions *consul.QueryOptions) (bool, *consul.TxnResponse, *consul.QueryMeta, error) {
		return true, nil, nil, nil
	}
	return m
}

func (m *MockClient) Txn(operations consul.TxnOps, queryOptions *consul.QueryOptions) (bool, *consul.TxnResponse, *consul.QueryMeta, error) {
	return m.TxnFunc(operations, queryOptions)
}
