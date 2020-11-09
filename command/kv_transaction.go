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
	"encoding/json"
	"fmt"

	"github.com/ExpediaGroup/flyte-client/flyte"
	client "github.com/ExpediaGroup/flyte-consul/client"
)

var (
	transactionSucceededEventDef  = flyte.EventDef{Name: "TransactionSucceeded"}
	transactionRolledBackEventDef = flyte.EventDef{Name: "TransactionRolledBack"}
)

//TransactKVInput represents the TransactKV command payload.
type TransactKVInput struct {
	Datacenter string               `json:"dc"`
	Operations []client.KVOperation `json:"operations"`
}

//TransactKVErrorOutput represents the error payload.
type TransactKVErrorOutput struct {
	Input  TransactKVInput             `json:"input"`
	Errors []client.KVTransactionError `json:"errors"`
}

//TransactKVResultOutput represents the result payload.
type TransactKVResultOutput struct {
	Input   TransactKVInput              `json:"input"`
	Results []client.KVTransactionResult `json:"results"`
}

//TransactKV produces the TransactKV flyte command.
func TransactKV(consulClient client.Consul) flyte.Command {
	return flyte.Command{
		Name: "TransactKV",
		OutputEvents: []flyte.EventDef{
			transactionSucceededEventDef,
			transactionRolledBackEventDef,
		},
		Handler: transactKVHandler(consulClient),
	}
}

func transactKVHandler(consulClient client.Consul) func(json.RawMessage) flyte.Event {
	return func(rawInput json.RawMessage) flyte.Event {
		input := TransactKVInput{}
		if err := json.Unmarshal(rawInput, &input); nil != err {
			return flyte.NewFatalEvent(fmt.Sprintf("input is not valid: %v", err))
		}
		if 0 == len(input.Operations) {
			return flyte.NewFatalEvent("missing operations")
		}

		errors := []client.KVTransactionError{}
		for index, operation := range input.Operations {
			isKeyMissing := operation.Key == ""
			isVerbSupported := consulClient.IsVerbSupported(operation.Verb)

			if isKeyMissing {
				errors = append(errors, client.KVTransactionError{
					Index: index,
					Error: "key is missing",
				})
			}
			if !isVerbSupported {
				errors = append(errors, client.KVTransactionError{
					Index: index,
					Error: fmt.Sprintf("%v verb is not valid", operation.Verb),
				})
			}
			if isKeyMissing || !isVerbSupported {
				continue
			}
		}

		if 0 != len(errors) {
			return newTransactionRolledBackEvent(input, errors)
		}

		results, rollback, err := consulClient.KVTransact(input.Datacenter, input.Operations)
		if nil != err {
			return flyte.NewFatalEvent(fmt.Sprintf("failed to make transaction request: %v", err))
		}
		if 0 < len(rollback) {
			return newTransactionRolledBackEvent(input, rollback)
		}

		return newTransactionSucceededEvent(input, results)
	}
}

func newTransactionRolledBackEvent(input TransactKVInput, errors []client.KVTransactionError) flyte.Event {
	return flyte.Event{
		EventDef: transactionRolledBackEventDef,
		Payload: TransactKVErrorOutput{
			Input:  input,
			Errors: errors,
		},
	}
}

func newTransactionSucceededEvent(input TransactKVInput, results []client.KVTransactionResult) flyte.Event {
	return flyte.Event{
		EventDef: transactionSucceededEventDef,
		Payload: TransactKVResultOutput{
			Input:   input,
			Results: results,
		},
	}
}
