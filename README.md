# Respond Collector

[![Build Status](https://travis-ci.org/FreifunkBremen/respond-collector.svg?branch=master)](https://travis-ci.org/FreifunkBremen/respond-collector)
[![Coverage Status](https://coveralls.io/repos/github/FreifunkBremen/respond-collector/badge.svg?branch=master)](https://coveralls.io/github/FreifunkBremen/respond-collector?branch=master)

`respond-collector` is a respond client that fetches, stores and publishes information about a Freifunk network. The goals:
* Generating JSON for [MeshViewer](https://github.com/ffnord/meshviewer) (Works with branch [JSONv2](https://github.com/FreifunkBremen/meshviewer/tree/JSONv2))
* Storing statistics in [InfluxDB](https://influxdata.com/) to be analyzed by [Grafana](http://grafana.org/)
* Provide information via JSON-APIs

## Usage
```
Usage of ./respond-collector:
  -config path/to/config.yml
```

## Development
### respond
It send the `gluon-neighbour-info` request and collect them together.

It will send UDP packetes by the multicast group `ff02:0:0:0:0:0:2:1001` and port `1001`.

### nodes.Nodes
It cached the information of the Nodes and will save them periodical to a JSON file.
The current nodes are saved default under `nodes.json`.

## Related projects

Collecting data from respondd:
* [Node informant](https://github.com/ffdo/node-informant) written in Go
* [HopGlass Server](https://github.com/plumpudding/hopglass-server) written in Node.js

Respondd for servers:
* [respondd branch of ffnord-alfred-announce](https://github.com/ffnord/ffnord-alfred-announce/tree/respondd) from FreiFunkNord
* [respondd](https://github.com/Sunz3r/ext-respondd) from Sunz3r
* [respondd](https://github.com/FreifunkBremen/respondd) from Freifunk Bremen (just a proof of concept)
