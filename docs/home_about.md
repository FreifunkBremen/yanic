# About

A little overview of yanic in connection with other software:
![Overview](overview.svg)

## How respondd works

It sends the `gluon-neighbour-info` request and collects the answers.

It will send UDP packets with multicast address `ff02:0:0:0:0:0:2:1001` and port `1001`.

If a node does not answer, it will request with the last know address under the port `1001`.

## Related projects

#### yanic collecting data
VPNs (respondd for servers):

* [mesh-announce](https://github.com/ffnord/mesh-announce) from FreiFunkNord
* [respondd](https://github.com/Sunz3r/ext-respondd) from Sunz3r

Nodes (respondd for nodes): [gluon](https://github.com/freifunk-gluon/gluon/)

#### Alternative collectors of respondd data:

* [Node informant](https://github.com/ffdo/node-informant) written in Go
* [HopGlass Server](https://github.com/plumpudding/hopglass-server) written in Node.js

#### yanic published data

**Databases:**

* [InfluxDB](https://influxdata.com/) SQL-like timeserial database
* [Graphite](https://graphiteapp.org/) RRD file Based

	Visualization from Databases: [Grafana](https://grafana.com/)

**Output:**
* meshviewer-ffrgb:
  * [meshviewer](https://github.com/ffrgb/meshviewer)
* nodelist:
  * [ffapi](https://freifunk.net/api-generator/)
    * [freifunk-karte.de](https://freifunk-karte.de)
* meshviewer (others):
  *  unmaintained [origin meshviewer](https://github.com/ffnord/meshviewer) branch: master (v1) and dev (v2)
