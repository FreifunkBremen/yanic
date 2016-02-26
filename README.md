# README
`micro-daemon` is a respond daemon to fetch information for Freifunk Nodes
and represent this information with Websocket- and JSON-APIs.

## Issues/Features in the Future
It will also APIs for manipulate the fetched data
 and give a access for ansible to push changes to the nodes.

Also it's will push statistic informations to a influxdb.

## Usage
```
Usage of ./micro-daemon:
  -aliases string
    	path aliases.json file (default "webroot/aliases.json")
  -collectInterval int
    	interval for data collections (default 15)
  -h string
    	path aliases.json file
  -output string
    	path nodes.json file (default "webroot/nodes.json")
  -p string
    	path aliases.json file (default "8080")
  -saveInterval int
    	interval for data saving (default 60)
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

A Implementation of the connection to influxdb are also needed, maybe log a little bit to `telegraf` from influxdb.
