# Transmission Exporter for Prometheus [![Build Status](https://drone.github.matthiasloibl.com/api/badges/metalmatze/transmission-exporter/status.svg)](https://drone.github.matthiasloibl.com/metalmatze/transmission-exporter)

Prometheus exporter for [Transmission](https://transmissionbt.com/) metrics, written in Go.

### Installation

    $ go get github.com/metalmatze/transmission-exporter

### Build

    make

For development we encourage you to use `make install` instead, it's faster. 

### Docker

    docker pull metalmatze/transmission-exporter
    docker run -d -p 19091:19091 metalmatze/transmission-exporter

### Original authors of the Transmission package  
Tobias Blom (https://github.com/tubbebubbe/transmission)  
Long Nguyen (https://github.com/longnguyen11288/go-transmission)


