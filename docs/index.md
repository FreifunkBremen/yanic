# Home


    __   __          _
    \ \ / /_ _ _ __ (_) ___
     \ V / _` | '_ \| |/ __|
      | | (_| | | | | | (__
      |_|\__,_|_| |_|_|\___|
    Yet another node info collector

[![status-badge](https://ci.codeberg.org/api/badges/14886/status.svg)](https://ci.codeberg.org/repos/14886/branches/main)
[![codecov](https://codecov.io/gh/FreifunkBremen/yanic/branch/main/graph/badge.svg)](https://codecov.io/gh/FreifunkBremen/yanic)
[![Go Report Card](https://goreportcard.com/badge/github.com/FreifunkBremen/yanic)](https://goreportcard.com/report/github.com/FreifunkBremen/yanic)

`yanic` is a respondd client that fetches, stores and publishes information about a Freifunk network.

## The goals:

* Generating JSON for [Meshviewer](https://github.com/ffrgb/meshviewer)
* Storing statistics in [InfluxDB](https://influxdata.com/) or [Graphite](https://graphiteapp.org/) to be analyzed by [Grafana](http://grafana.org/)
* Provide a little webserver for a standalone installation with a meshviewer
