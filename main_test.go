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

package main

import (
	"testing"

	client "github.com/ExpediaGroup/flyte-consul/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPackDefinitionIsPopulated(t *testing.T) {
	packDef := GetPackDef(DummyConsul{})

	assert.Equal(t, "Consul", packDef.Name)
	assert.Equal(t, "https://github.com/ExpediaGroup/flyte-consul/blob/master/README.md", packDef.HelpURL.String())
	require.Equal(t, 0, len(packDef.Labels))
	require.Equal(t, 1, len(packDef.Commands))
	require.Equal(t, 0, len(packDef.EventDefs))
}

type DummyConsul struct{}

func (DummyConsul) KVTransact(datacenter string, operations []client.KVOperation) ([]client.KVTransactionResult, []client.KVTransactionError, error) {
	return nil, nil, nil
}

func (DummyConsul) IsVerbSupported(verb string) bool {
	return true
}
