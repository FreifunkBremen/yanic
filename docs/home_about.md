---
layout: default
title: About
permalink: /home/about.html
---

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
