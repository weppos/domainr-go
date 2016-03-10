# Domainr API client

A Go client for the [Domainr API](http://domainr.build/).

[![Build Status](https://travis-ci.org/weppos/domainr-go.svg?branch=master)](https://travis-ci.org/weppos/domainr-go)
[![GoDoc](https://godoc.org/github.com/weppos/domainr-go/domainr?status.svg)](https://godoc.org/github.com/weppos/domainr-go/domainr)


## Getting started

```shell
$ git clone git@github.com:weppos/domainr-go.git
$ cd domainr-go
```

Run the test suite.


## Testing

```shell
$ go test ./...
```

### Live Testing

```shell
$ export DOMAINR_CLIENT_ID="some-magic-client-id"
$ go test ./... -v
```

**Example output**

```shell
$ go test ./... -v
=== RUN   TestNewClient
--- PASS: TestNewClient (0.00s)
=== RUN   TestLivePrivateGetStatus
<nil>
&{[{domainr.com com active active}] 0xc820332b40}
--- PASS: TestLivePrivateGetStatus (1.13s)
=== RUN   TestLiveGetStatus
<nil>
&{domainr.com com active active}
--- PASS: TestLiveGetStatus (0.24s)
PASS
ok  	github.com/weppos/domainr-go/domainr	1.385s
```

**Custom domain list**

```shell
$ DOMAINR_STATUS_DOMAINS=dnsimple.com,domainr.com go test ./... -v
=== RUN   TestNewClient
--- PASS: TestNewClient (0.00s)
=== RUN   TestLivePrivateGetStatus
<nil>
&{[{dnsimple.com com active registrar registrar} {domainr.com com active active}] 0xc82041a090}
--- PASS: TestLivePrivateGetStatus (0.50s)
=== RUN   TestLiveGetStatus
<nil>
&{dnsimple.com com active registrar registrar}
--- PASS: TestLiveGetStatus (0.26s)
PASS
ok  	github.com/weppos/domainr-go/domainr	0.772s
```

## Installation

```shell
$ go get github.com/weppos/domainr-go/domainr
```

## Usage

```go
package main

import (
    "fmt"
    "os"

    "github.com/weppos/domainr-go/domainr"
)

func main() {
  clientID := "some-magic-client-id"

  client := domainr.NewClient(NewDomainrAuthentication(clientID))

  // Get the status of some domains
  domainResponse, err := client.Status([]string{"example.com", "example.org"})
  if err != nil {
      fmt.Println(err)
      os.Exit(1)
  }

  domain := domainResponse.Domains[0]
  fmt.Printf("%s: %s", domain.Name, domain.Summary)
}
```

### Authentication

This library supports both Mashape and Commercial authentication.

```go
// commercial authentication
client := domainr.NewClient(NewDomainrAuthentication("client-id"))

// mashape authentication
client := domainr.NewClient(NewMashapeAuthentication("mashape-api-key"))
```

The API endpoint is automatically adapted according to the type of authentication selected.


## License

Copyright (c) 2016 Simone Carletti. This is Free Software distributed under the MIT license.
