# Respond Collector

[![Build Status](https://travis-ci.org/FreifunkBremen/respond-collector.svg?branch=master)](https://travis-ci.org/FreifunkBremen/respond-collector)
[![Coverage Status](https://coveralls.io/repos/github/FreifunkBremen/respond-collector/badge.svg?branch=master)](https://coveralls.io/github/FreifunkBremen/respond-collector?branch=master)

`respond-collector` is a respondd client that fetches, stores and publishes information about a Freifunk network. The goals:
* Generating JSON for [Meshviewer](https://github.com/ffrgb/meshviewer)
* Storing statistics in [InfluxDB](https://influxdata.com/) to be analyzed by [Grafana](http://grafana.org/)
* Provide a little webserver for a standalone installation with a meshviewer

## Usage
```
Usage of ./respond-collector:
  -config path/to/config.toml
```
## Configuration
Read comments in [config_example.toml](config_example.toml) for more information.

## Live
* [meshviewer](https://map.bremen.freifunk.net) **Freifunk Bremen** with a patch to show state-version of `nodes.json`
* [grafana](https://grafana.bremen.freifunk.net)  **Freifunk Bremen** show data of InfluxDB

## How it works

It sends the `gluon-neighbour-info` request and collects the answers.

It will send UDP packets with multicast group `ff02:0:0:0:0:0:2:1001` and port `1001`.

If a node does not answer, it will request with the last know address under the port `1001`.


## Related projects

Collecting data from respondd:
* [Node informant](https://github.com/ffdo/node-informant) written in Go
* [HopGlass Server](https://github.com/plumpudding/hopglass-server) written in Node.js

Respondd for servers:
* [ffnord-alfred-announce](https://github.com/ffnord/ffnord-alfred-announce) from FreiFunkNord
* [respondd](https://github.com/Sunz3r/ext-respondd) from Sunz3r
