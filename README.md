[![Build Status](https://travis-ci.org/platanus/lightning-prometheus-exporter.svg?branch=master)](https://travis-ci.org/platanus/lightning-prometheus-exporter) [![](https://images.microbadger.com/badges/version/platanus/lightning-prometheus-exporter.svg)](https:/hub.docker.com/r/platanus/lightning-prometheus-exporter) [![Go Report Card](https://goreportcard.com/badge/github.com/platanus/lightning-prometheus-exporter)](https://goreportcard.com/report/github.com/platanus/lightning-prometheus-exporter)

# lightning Prometheus Exporter

lightning Prometheus exporter makes it possible to monitor lightning node using Prometheus.

## Overview

Lightning Prometheus exporter fetches the metrics the rpc api of the node, converts the metrics into appropriate Prometheus metrics types and finally exposes them via an HTTP server to be collected by [Prometheus](https://prometheus.io/).

## Getting Started

In this section, we show how to quickly run lightning Prometheus Exporter for lightning.

### Prerequisites

We assume that you have already installed Prometheus and a LND Lightning Node. Additionally, you need to configure Prometheus to scrape metrics from the server with the exporter. Note that the default scrape port of the exporter is `9113` and the default metrics path -- `/metrics`.

## Usage

### Command-line Arguments

```
Usage of ./lightning-prometheus-exporter:
  -namepsace string
        The namespace or prefix to use in the exported metrics. The default value can be overwritten by NAMESPACE environment variable (default: lnd)
  -web.telemetry-path string
        A path under which to expose metrics. The default value can be overwritten by TELEMETRY_PATH environment variable. (default "/metrics")
  -web.listen-address string
        An address to listen on for web interface and telemetry. The default value can be overwritten by LISTEN_ADDRESS environment variable. (default ":9113")
  -rpc.host string
        Lightning node RPC host. The default value can be overwritten by RPC_HOST environment variable (default: localhost)
  -rpc.Port int
        Lightning node RPC port. The default value can be overwritten by RPC_PORT environment variable (default: 10009)
  -lnd.tls-cert-path string
        The path to the tls certificate. The default value can be overwritten by TLS_CERT_PATH environment variable (default: "/root/.lnd")
  -lnd.macaroon-path
        The path to the read only macaroon. The default value can be overwritten by MACAROON_PATH environment variable
```

### Exported Metrics

* Connect to the `/metrics` page of the running exporter to see the complete list of metrics along with their descriptions.

### Troubleshooting

The exporter logs errors to the standard output. When using Docker, if the exporter doesn’t work as expected, check its logs using [docker logs](https://docs.docker.com/engine/reference/commandline/logs/) command.

## Releases

For each release, we publish the corresponding Docker image at `platanus/lightning-prometheus-exporter` [DockerHub repo](https://hub.docker.com/r/platanus/lightning-prometheus-exporter/) and the binaries on the GitHub [releases page](https://github.com/platanus/lightning-prometheus-exporter/releases).

## Building the Exporter

You can build the exporter using the provided Makefile. Before building the exporter, make sure the following software is installed on your machine:
* make
* git
* Docker for building the container image
* Go for building the binary

### Building the Docker Image

To build the Docker image with the exporter, run:
```
$ make container
```

Note: go is not required, as the exporter binary is built in a Docker container. See the [Dockerfile](Dockerfile).

### Building the Binary

To build the binary, run:
```
$ make
```

Note: the binary is built for the OS/arch of your machine. To build binaries for other platforms, see the [Makefile](Makefile).

The binary is built with the name `lightning-prometheus-exporter`.

## Credits

Thank you [contributors](https://github.com/platanus/lightning-prometheus-exporter/graphs/contributors)!

<img src="http://platan.us/gravatar_with_text.png" alt="Platanus" width="250"/>

lightning Prometheus Exporter is maintained by [platanus](http://platan.us).

## License

Lightning Prometheus Exporter is © 2018 platanus, spa. It is free software and may be redistributed under the terms specified in the LICENSE file.
