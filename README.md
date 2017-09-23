# drone-crowdin

[![Go Doc](https://godoc.org/github.com/JonasFranzDEV/drone-crowdin?status.svg)](http://godoc.org/github.com//JonasFranzDEV/drone-crowdin)
[![Go Report](https://goreportcard.com/badge/github.com/JonasFranzDEV/drone-crowdin)](https://goreportcard.com/report/github.com//JonasFranzDEV/drone-crowdin)

Drone plugin to update translation files at Crowdin. For the usage information and a listing of the available options please take a look at [the docs](DOCS.md).

## Build

Build the binary with the following commands:

```
go build
```

## Docker

Build the Docker image with the following commands:

```
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -tags netgo -o release/linux/amd64/drone-crowdin
docker build --rm -t plugins/crowdin .
```
