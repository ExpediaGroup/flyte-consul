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
	"net/url"
	"os"

	"github.com/HotelsDotCom/go-logger"
)

const (
	flyteHostKey = "FLYTE_API"
	packNameKey  = "PACK_NAME"
)

var lookupEnv = os.LookupEnv

func flyteAPIHost() *url.URL {
	hostEnv := getEnv(flyteHostKey, true)
	host, err := url.Parse(hostEnv)
	if err != nil {
		logger.Fatalf("%s=%p is not a valid URL: %v", flyteHostKey, hostEnv, err)
	}
	return host
}

func packName() string {
	return getEnv(packNameKey, false)
}

func getEnv(key string, required bool) string {
	if v, _ := lookupEnv(key); v != "" {
		return v
	}

	if required {
		logger.Fatalf("%s env not set", key)
	}

	return ""
}
