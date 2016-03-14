# Respond Collector

`Respond Collector` is a respond client that fetches, stores and publishes information about a Freifunk network. The goals:
* Generating JSON for [MeshViewer](https://github.com/ffnord/meshviewer) (Works with branch [JSONv2](https://github.com/FreifunkBremen/meshviewer/tree/JSONv2))
* Storing statistics in [InfluxDB](https://influxdata.com/) to be analyzed by [Grafana](http://grafana.org/)
* Provide information via Websocket- and JSON-APIs

## Usage
```
Usage of ./RespondCollector:
  -config path/to/config.yml
```

## Development
### respond
It send the `gluon-neighbour-info` request and collect them together.

It will send UDP packetes by the multicast group `ff02:0:0:0:0:0:2:1001` and port `1001`.

### modes.Nodes
It cached the Informations of the Nodes and will save them periodical to a JSON file.
The current nodes are saved default under `nodes.json`.


### websocketserver
One Instance is running under `/nodes` which send updates or new Nodes,
 which are collected by respond.

### Issues
Later there should be also `/aliases` Websocket with Authentification to manage the `aliases.json` with the request for changes.

