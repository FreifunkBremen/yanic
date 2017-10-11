---
layout: default
title: Home
---

    __   __          _
    \ \ / /_ _ _ __ (_) ___
     \ V / _` | '_ \| |/ __|
      | | (_| | | | | | (__
      |_|\__,_|_| |_|_|\___|
    Yet another node info collector

\(previously [respond-collector](https://github.com/FreifunkBremen/respond-collector)\)

[![Build Status](https://travis-ci.org/FreifunkBremen/yanic.svg?branch=master)](https://travis-ci.org/FreifunkBremen/yanic) [![](https://coveralls.io/repos/github/FreifunkBremen/yanic/badge.svg?branch=master)](https://coveralls.io/github/FreifunkBremen/yanic?branch=master) [![](https://goreportcard.com/badge/github.com/FreifunkBremen/yanic)](https://goreportcard.com/report/github.com/FreifunkBremen/yanic)

`yanic` is a respondd client that fetches, stores and publishes information about a Freifunk network.

## The goals:

* Generating JSON for [Meshviewer](https://github.com/ffrgb/meshviewer)
* Storing statistics in [InfluxDB](https://influxdata.com/) or [Graphite](https://graphiteapp.org/) to be analyzed by [Grafana](http://grafana.org/)
* Provide a little webserver for a standalone installation with a meshviewer
