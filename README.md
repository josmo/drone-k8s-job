# drone-k8s-job

[![Build Status](https://drone.pelo.tech/api/badges/josmo/drone-k8s-job/status.svg)](https://drone.pelo.tech/josmo/drone-k8s-job)
[![Go Doc](https://godoc.org/github.com/josmo/drone-k8s-job?status.svg)](http://godoc.org/github.com/josmo/drone-k8s-job)
[![Go Report](https://goreportcard.com/badge/github.com/josmo/drone-k8s-job)](https://goreportcard.com/report/github.com/josmo/drone-k8s-job)
[![](https://images.microbadger.com/badges/image/peloton/drone-k8s-job.svg)](https://microbadger.com/images/peloton/drone-k8s-job "Get your own image badge on microbadger.com")

Drone plugin to deploy a job in k8s. For the usage information and a listing of the available options please take a look at [the docs](DOCS.md).
 
## Experimental!!!!

This is currently in an experimental phase. Currently insecure and cert aren't implemented. Please feel free to provide feedback and suggestions


## Versions

This repo is using auto-tag from the drone-docker plugin meaning that
1. master will always publish to 'latest' in docker hub peloton/drone-k8s-job
2. tags will follow semver at the 1.0.0+ - initial 0.x.x may have breaking changes

## Binary

Build the binary using `go build`:


### Example

## Usage

Build and deploy from your current working directory:

```
docker run --rm                          \
  -e PLUGIN_URL=<source>                 \
  -e PLUGIN_TOKEN=<token>                \
  -e PLUGIN_CERT=<cert>                  \
  -e PLUGIN_INSECURE=<true>              \
  -e PLUGIN_NAMESPACES=<namespaces>      \
  -e JOB_TEMPLATE=job.yml                |
  -v $(pwd):$(pwd)                       \
  -w $(pwd)                              \
  peloton/drone-k8s-job 
```
