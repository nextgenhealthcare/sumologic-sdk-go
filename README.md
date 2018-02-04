[![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/brandonstevens/sumologic-sdk-go)
[![license](https://img.shields.io/github/license/mashape/apistatus.svg)](https://choosealicense.com/licenses/mit/)
[![Build Status](https://travis-ci.org/brandonstevens/sumologic-sdk-go.svg)](https://travis-ci.org/brandonstevens/sumologic-sdk-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/brandonstevens/sumologic-sdk-go)](https://goreportcard.com/report/github.com/brandonstevens/sumologic-sdk-go)

# SumoLogic API in Go

A Go wrapper for the SumoLogic API.

## Contributing

Any and all contributions are welcome. Please don't hestiate to submit an issue or pull request.

## Roadmap

The initial release is focused on being consumed by a Terraform provider in AWS environments such as support for managing hosted collectors and AWS specific hosted sources (e.g. AWS Cloudtrail).

## Installation

```go
import "github.com/brandonstevens/sumologic-sdk-go"
```

## Usage

* auth_token: Base64 encoding of `<accessId>:<accessKey>`. For more information, see [API Authentication](https://help.sumologic.com/APIs/General-API-Information/API-Authentication)
* endpoint_url: Sumo Logic has several deployments that are assigned depending on the geographic location and the date an account is created. For more information, see [Sumo Logic Endpoints and Firewall Security](https://help.sumologic.com/APIs/General-API-Information/Sumo-Logic-Endpoints-and-Firewall-Security)

```go
client, _ := sumologic.NewClient("auth_token", "endpoint_url")

collector, _, err := client.GetHostedCollector(134485191)
if err == sumologic.ErrCollectorNotFound {
	log.Fatalf("Collector not found: %s\n", err)
}
if err != nil {
	log.Fatalf("Unknown error: %s\n", err)
}

log.Printf("Collector %d: %s\n", collector.Id, collector.Name)
```

## Development

Run unit tests with `make test`.