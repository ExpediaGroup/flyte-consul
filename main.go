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
	"log"
	"net/url"
	"os"
	"time"

	flyteClient "github.com/ExpediaGroup/flyte-client/client"
	"github.com/ExpediaGroup/flyte-client/flyte"
	client "github.com/ExpediaGroup/flyte-consul/client"
	"github.com/ExpediaGroup/flyte-consul/command"
	"github.com/HotelsDotCom/go-logger"
)

const (
	packDefHelpURL  = "https://github.com/ExpediaGroup/flyte-consul/blob/master/README.md"
	defaultPackName = "Consul"
)

func main() {
	consulClient, err := client.NewConsul()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	packDef := GetPackDef(consulClient)
	pack := flyte.NewPack(packDef, flyteClient.NewClient(flyteAPIHost(), 10*time.Second))
	pack.Start()

	select {}
}

// GetPackDef gets the flight pack definition.
func GetPackDef(consul client.Consul) flyte.PackDef {
	helpURL, err := url.Parse(packDefHelpURL)
	if err != nil {
		logger.Fatal("invalid pack help url")
	}

	packName := packName()
	if packName == "" {
		packName = defaultPackName
	}

	return flyte.PackDef{
		Name:    packName,
		HelpURL: helpURL,
		Commands: []flyte.Command{
			command.TransactKV(consul),
		},
		EventDefs: []flyte.EventDef{},
	}
}
