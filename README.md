# Yanic
```
__   __          _
\ \ / /_ _ _ __ (_) ___
 \ V / _` | '_ \| |/ __|
  | | (_| | | | | | (__
  |_|\__,_|_| |_|_|\___|
Yet another node info collector
```

[![Build Status](https://travis-ci.org/FreifunkBremen/yanic.svg?branch=master)](https://travis-ci.org/FreifunkBremen/yanic)
[![Coverage Status](https://coveralls.io/repos/github/FreifunkBremen/yanic/badge.svg?branch=master)](https://coveralls.io/github/FreifunkBremen/yanic?branch=master)

`yanic` is a respondd client that fetches, stores and publishes information about a Freifunk network. The goals:
* Generating JSON for [Meshviewer](https://github.com/ffrgb/meshviewer)
* Storing statistics in [InfluxDB](https://influxdata.com/) to be analyzed by [Grafana](http://grafana.org/)
* Provide a little webserver for a standalone installation with a meshviewer

## Usage
```
Usage of ./yanic:
  -config path/to/config.toml
```
## Configuration
Read comments in [config_example.toml](config_example.toml) for more information.

## Live
* [meshviewer](https://map.bremen.freifunk.net) **Freifunk Bremen** with a patch to show state-version of `nodes.json`
* [grafana](https://grafana.bremen.freifunk.net)  **Freifunk Bremen** show data of InfluxDB

## How it works

In the first step Yanic sends a multicast message to the group `ff02:0:0:0:0:0:2:1001` and port `1001`.
Recently seen nodes that does not reply are requested via a unicast message.

## Related projects

Collecting data from respondd:
* [HopGlass Server](https://github.com/plumpudding/hopglass-server) written in Node.js

Respondd for servers:
* [ffnord-alfred-announce](https://github.com/ffnord/ffnord-alfred-announce) from FreiFunkNord
* [respondd](https://github.com/Sunz3r/ext-respondd) from Sunz3r
