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
	"os"
	"testing"

	"github.com/HotelsDotCom/go-logger/loggertest"
	"github.com/stretchr/testify/assert"
)

var TestEnv map[string]string

func BeforeConfig() {
	loggertest.Init("DEBUG")
	TestEnv = make(map[string]string)
	lookupEnv = func(key string) (string, bool) {
		value, ok := TestEnv[key]
		return value, ok
	}
}

func AfterConfig() {
	lookupEnv = os.LookupEnv
	loggertest.Reset()
}

func TestFlyteAPIHostEnv(t *testing.T) {
	BeforeConfig()
	defer AfterConfig()

	TestEnv["FLYTE_API"] = "http://joe.mama.com:8080"

	url := flyteAPIHost()
	assert.Equal(t, "http://joe.mama.com:8080", url.String())
}

func TestFlyteAPIHostEnvNotSet(t *testing.T) {
	BeforeConfig()
	defer AfterConfig()

	assert.Panics(t, func() { flyteAPIHost() })
}

func TestFlyteAPIHostEnvInvalidURL(t *testing.T) {
	BeforeConfig()
	defer AfterConfig()

	TestEnv["FLYTE_API"] = ":invalid. url:BOBO"

	assert.Panics(t, func() { flyteAPIHost() })
}

func TestPackName(t *testing.T) {
	BeforeConfig()
	defer AfterConfig()

	TestEnv["PACK_NAME"] = "Consul2"
	assert.Equal(t, "Consul2", packName())
}
