# flyte-consul

![Build Status](https://travis-ci.org/ExpediaGroup/flyte-consul.svg?branch=master)
[![Docker Stars](https://img.shields.io/docker/stars/expediagroup/flyte-consul.svg)](https://hub.docker.com/r/expediagroup/flyte-consul)
[![Docker Pulls](https://img.shields.io/docker/pulls/expediagroup/flyte-consul.svg)](https://hub.docker.com/r/expediagroup/flyte-consul)

A Consul pack for flyte.

## Build

Pack requires go version min. 1.14

- go build `go build`
- go test `go test ./...`
- docker build `docker build -t <name>:<version> .`

## Configuration

The plugin is configured using environment variables:

ENV VAR                          | Default  |  Description                               | Example               
 ------------------------------- |  ------- |  ----------------------------------------- |  ---------------------
FLYTE_API                        | -        | The API endpoint to use                    | http://localhost:8080

See [consul documentation](https://www.consul.io/commands#environment-variables) for consul specific environment variables.

Example `FLYTE_API=http://localhost:8080 ./flyte-consul`

## Commands

### TransactKV

    {
        "dc": "...", // optional
        "operations": [ // required (at least one)
            {
                "verb": "...", // required (see https://godoc.org/github.com/hashicorp/consul/api#KVOp for supported values)
                "key": "...", // required
                "value": {...} // optional
            },
            ...
        ]
    }

#### Returned events

`TransactionSucceeded`

    {
        "input": {
            "dc": "...",
            "operations": [
                {
                    "verb": "...",
                    "key": "...",
                    "value": {...}
                },
                ...
            ]
        },
        "results": [
            {
                "index": 0..n,
                "key": "...",
                "value": {...}
            },
            ...
        ]
    }

`TransactionRolledBack`

    {
        "input": {
            "dc": "...",
            "operations": [
                {
                    "verb": "...",
                    "key": "...",
                    "value": {...}
                },
                ...
            ]
        },
        "errors": [
            {
                "index": 0..n,
                "error": "..."
            },
            ...
        ]
    }

# consul-flyte-pack

## Prerequisites

Building the application requires [Go](https://golang.org/). Please install Go as instructed below
or using [Docker](https://www.docker.com/) to compile and run the application.

The application is using [Modules](https://github.com/golang/go/wiki/Modules) for dependency management.

### Installing Go tools

Instructions: https://golang.org/doc/install.

TLDR: download Go binaries and add them to `PATH` environment variable.

```
curl -L 'https://dl.google.com/go/go1.11.2.darwin-amd64.tar.gz' | tar xz -C "$HOME/Downloads"
export PATH="$PATH:$HOME/Downloads/go/bin"
# Test it
go version
```

If you are using Homebrew you can also install Go by running:

```
brew install go
```

By default `$GOPATH` is set to `$HOME/go` so all packages will be downloaded into this folder.

For release build, it is recommended to use `./go.sh` wrapper script as build information such as
git commit, branch and build time will be included in the binary.
The script can also download Go into `$GOPATH` when Go is not installed and use it to run Go commands.

### Installing Docker

If you want to use Docker, please follow the instructions at https://docs.docker.com/engine/installation/

The `./build-auto.sh` script utilizes Docker and is used to build and test the application on
automation server like Jenkins.

## Development

### Build and test

Download dependencies:

```
./go.sh mod download
```

Build:

```
./go.sh build -v
```

Run tests and benchmarks:

```
./go.sh test -v -bench . ./...
```

Run the application:

```
APP_NAME=consul-flyte-pack CERT_PASS=changeit ./consul-flyte-pack
```

Open browser and try https://localhost:8443/buildInfo

### Docker

Run tests in Docker container:

```
./go.sh docker go test -v ./...
```

Build locally and run the binary in Docker container:

```
GOOS=linux GOARCH=amd64 ./go.sh build
docker build -t consul-flyte-pack:local .
docker run -e "APP_NAME=consul-flyte-pack" -p 8443:8443 consul-flyte-pack:local
```

### Vault

In case you get vault "certificate signed by unknown authority" error while developing from your machine, set this environment variable:

```
VAULT_CACERT=/path/to/consul-flyte-pack/conf/vault-ca.pem
```

`VAULT_SSL_CERT` is also supported for compability purpose.

### Learning Go 

Writing code: https://golang.org/doc/code.html

Editors: https://golang.org/doc/editors.html

For more information, see https://golang.org/doc/

## Troubleshooting

### Could not download packages from GitHub
```
go: finding github.com/a/package vx.y.z
(exit code 128)
```

It can be GitHub throttled the requests.
Try ssh protocol instead: `git config --global url.git@github.com:.insteadOf https://github.com/`

## Legal
This project is available under the [Apache 2.0 License](http://www.apache.org/licenses/LICENSE-2.0.html).

Copyright 2020 Expedia, Inc.
