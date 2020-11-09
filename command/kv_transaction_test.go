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

package command

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ExpediaGroup/flyte-consul/client"
	"github.com/HotelsDotCom/go-logger/loggertest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var KVTransactionMockConsul *MockConsul

func Before() {
	KVTransactionMockConsul = NewMockConsul()
	loggertest.Init("DEBUG")
}

func After() {
	loggertest.Reset()
}

func TestTransactKVCommandIsPopulated(t *testing.T) {
	Before()
	defer After()

	command := TransactKV(KVTransactionMockConsul)

	assert.Equal(t, "TransactKV", command.Name)
	require.Equal(t, 2, len(command.OutputEvents))
	assert.Equal(t, "TransactionSucceeded", command.OutputEvents[0].Name)
	assert.Equal(t, "TransactionRolledBack", command.OutputEvents[1].Name)
}

func TestTransactKVReturnsTransactionSucceededEvent(t *testing.T) {
	Before()
	defer After()

	const key = "joe/mama"
	KVTransactionMockConsul.KVTransactFunc = func(datacenter string, operations []client.KVOperation) ([]client.KVTransactionResult, []client.KVTransactionError, error) {
		return []client.KVTransactionResult{
			{
				Index: 0,
				Key:   key,
				Value: nil,
			},
		}, []client.KVTransactionError{}, nil
	}
	KVTransactionMockConsul.IsVerbSupportedFunc = func(verb string) bool {
		return true
	}

	handler := TransactKV(KVTransactionMockConsul).Handler

	event := handler(getValidTransactKVPayload())
	println(fmt.Sprintf("event: %+v", event.Payload))
	output := event.Payload.(TransactKVResultOutput)

	require.NotNil(t, event)
	assert.Equal(t, "TransactionSucceeded", event.EventDef.Name)
	assert.Equal(t, "dc", output.Input.Datacenter)
	require.Equal(t, 1, len(output.Input.Operations))
	assert.Equal(t, key, output.Input.Operations[0].Key)
	assert.Less(t, 0, len(output.Input.Operations[0].Value))
	assert.Equal(t, "set", output.Input.Operations[0].Verb)
}

func TestTransactKVFailsInvalidInput(t *testing.T) {
	Before()
	defer After()

	handler := TransactKV(KVTransactionMockConsul).Handler
	event := handler([]byte(`asdk292ds{}][;dsfjIljskdf{}[`))

	require.NotNil(t, event)
	assert.Equal(t, "FATAL", event.EventDef.Name)
	assert.True(t, strings.HasPrefix(event.Payload.(string), "input is not valid"))
}

func TestTransactKVFailsMissingOperations(t *testing.T) {
	Before()
	defer After()

	handler := TransactKV(KVTransactionMockConsul).Handler
	event := handler([]byte(`{
		"dc": "dc",
		"operations": []
	}`))

	require.NotNil(t, event)
	assert.Equal(t, "FATAL", event.EventDef.Name)
	assert.Equal(t, "missing operations", event.Payload)
}

func TestTransactKVFailedMissingKey(t *testing.T) {
	Before()
	defer After()

	KVTransactionMockConsul.IsVerbSupportedFunc = func(verb string) bool {
		return true
	}

	handler := TransactKV(KVTransactionMockConsul).Handler
	event := handler([]byte(`{
		"dc": "dc",
		"operations": [
			{
				"verb": "set",
				"value": {}
			}
		]
	}`))

	require.NotNil(t, event)
	assert.Equal(t, "TransactionRolledBack", event.EventDef.Name)
	output := event.Payload.(TransactKVErrorOutput)
	require.Equal(t, 1, len(output.Errors))
	assert.Equal(t, 0, output.Errors[0].Index)
	assert.Equal(t, "key is missing", output.Errors[0].Error)
}

func TestTransactKVFailedUnsupportedVerb(t *testing.T) {
	Before()
	defer After()

	KVTransactionMockConsul.IsVerbSupportedFunc = func(verb string) bool {
		return false
	}

	handler := TransactKV(KVTransactionMockConsul).Handler
	event := handler(getValidTransactKVPayload())

	require.NotNil(t, event)
	assert.Equal(t, "TransactionRolledBack", event.EventDef.Name)
	output := event.Payload.(TransactKVErrorOutput)
	require.Equal(t, 1, len(output.Errors))
	assert.Equal(t, 0, output.Errors[0].Index)
	assert.True(t, strings.HasSuffix(output.Errors[0].Error, "verb is not valid"))
}

func TestTransactKVRequestFailed(t *testing.T) {
	Before()
	defer After()

	KVTransactionMockConsul.KVTransactFunc = func(datacenter string, operations []client.KVOperation) ([]client.KVTransactionResult, []client.KVTransactionError, error) {
		return nil, nil, fmt.Errorf("kablammo")
	}
	KVTransactionMockConsul.IsVerbSupportedFunc = func(verb string) bool {
		return true
	}

	handler := TransactKV(KVTransactionMockConsul).Handler
	event := handler(getValidTransactKVPayload())

	require.NotNil(t, event)
	assert.Equal(t, "FATAL", event.EventDef.Name)
	assert.Equal(t, "failed to make transaction request: kablammo", event.Payload)
}

func TestTransactKVRolledBack(t *testing.T) {
	Before()
	defer After()

	const errorMsg = "kablammo"
	KVTransactionMockConsul.IsVerbSupportedFunc = func(verb string) bool {
		return true
	}
	KVTransactionMockConsul.KVTransactFunc = func(datacenter string, operations []client.KVOperation) ([]client.KVTransactionResult, []client.KVTransactionError, error) {
		return nil, []client.KVTransactionError{
			{
				Index: 0,
				Error: errorMsg,
			},
		}, nil
	}

	handler := TransactKV(KVTransactionMockConsul).Handler
	event := handler(getValidTransactKVPayload())

	require.NotNil(t, event)
	assert.Equal(t, "TransactionRolledBack", event.EventDef.Name)
	output := event.Payload.(TransactKVErrorOutput)
	require.Equal(t, 1, len(output.Errors))
	assert.Equal(t, 0, output.Errors[0].Index)
	assert.Equal(t, errorMsg, output.Errors[0].Error)
}

func getValidTransactKVPayload() []byte {
	return []byte(getValidTransactKVPayloadString())
}

func getValidTransactKVPayloadString() string {
	return `{
		"dc": "dc",
		"operations": [
			{
				"verb": "set",
				"key": "joe/mama",
				"value": {
					"one": 1,
					"two": "2",
					"three": true,
					"four": null,
					"five": []
				}
			}
		]
	}`
}

func NewMockConsul() *MockConsul {
	m := &MockConsul{}
	return m
}

type MockConsul struct {
	KVTransactFunc      func(datacenter string, operations []client.KVOperation) ([]client.KVTransactionResult, []client.KVTransactionError, error)
	IsVerbSupportedFunc func(verb string) bool
}

func (m *MockConsul) KVTransact(datacenter string, operations []client.KVOperation) ([]client.KVTransactionResult, []client.KVTransactionError, error) {
	return m.KVTransactFunc(datacenter, operations)
}

func (m *MockConsul) IsVerbSupported(verb string) bool {
	return m.IsVerbSupportedFunc(verb)
}
