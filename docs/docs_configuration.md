---
layout: default
title: Configuration
category_id: 1
permalink: /docs/configuration.html
---

Here you would find a long description, maybe the description in [example file](https://github.com/FreifunkBremen/yanic/blob/master/config_example.toml) are enough for you.

* TOC
{:toc}

## [respondd]
Group for configuration of respondd request.

```toml
[respondd]
enable           = true
#synchronize      = "1m"
collect_interval = "1m"
interface        = "eth0"
#port            = 10001
```


### enable
Enable request and collection of data per respondd requests

```toml
enable = true
```



### synchronize
Delay startup until a multiple of the period since zero time

```toml
synchronize      = "1m"
```



### collect_interval
How oftern send request per respondd.

It will send UDP packets with multicast group `ff02::2:1001` and port `1001`.
If a node does not answer after the half time, it will request with the last know address under the port `1001`.

```toml
collect_interval = "1m"
```



### interface
On which interface listen and send request.

```toml
interface        = "eth0"
```


### port
Define a port to listen and send the respondd packages.
If it is not set or 0 would choose the kernel a free port at his own

```toml
port              = 10001
```



## [webserver]
Yanic has a little build-in webserver, which statically serves a directory.
This is useful for testing purposes or for a little standalone installation.

```toml
[webserver]
enable  = false
bind    = "127.0.0.1:8080"
webroot = "/var/www/html/meshviewer"
```


### enable
Enable to start the built-in webserver of Yanic

```toml
enable  = false
```


### bind
On which ip address and port listen the webserver
```toml
bind    = "127.0.0.1:8080"
```


### webroot
The path to a folder, which files are published on this webserver.

```toml
webroot = "/var/www/html/meshviewer"
```



## [nodes]
```toml
[nodes]
enable         = true
state_path     = "/var/lib/yanic/state.json"
save_interval = "5s"
offline_after = "10m"
prune_after = "7d"
```


### enable
Enable the storing and writing of json files

```toml
enable         = true
```


### state_path
State-version of nodes.json to store cached data, these is the directly collected respondd data.

```toml
state_path     = "/var/lib/yanic/state.json"
```


### save_interval
Export nodes and graph periodically.

```toml
save_interval = "5s"
```


### offline_after
Set node to offline if not seen within this period.

```toml
offline_after = "10m"
```


### prune_after
Prune offline nodes after a time of inactivity.

```toml
prune_after = "7d"
```


## [meshviewer]
```toml
[meshviewer]
version        = 2
nodes_path     = "/var/www/html/meshviewer/data/nodes.json"
graph_path     = "/var/www/html/meshviewer/data/graph.json"
```

### version
The structur of nodes.json of which `nodes_path` `nodes.json` should saved:
* version 1 is to support legacy meshviewer (which are in master branch)
* https://github.com/ffnord/meshviewer/tree/master
* version 2 is to support new version of meshviewer (which are in legacy develop branch or newer)
* https://github.com/ffnord/meshviewer/tree/dev
* https://github.com/ffrgb/meshviewer/tree/develop

```toml
version = 2
```


### nodes_path
The path, where to store nodes.json (supports version 1 and two, see `nodes_version`)

```toml
nodes_path = "/var/www/html/meshviewer/data/nodes.json"
```


### graph_path
The path, where to store graph.json (only version 1)

```toml
graph_path = "/var/www/html/meshviewer/data/graph.json"
```



## [database]
The database organize all database types.
For all database types the is a internal job, which reset data for nodes (global statistics are still stored).
_(We have for privacy policy to store node data for maximum seven days.)_

Every database type has his own configuration under `database.connection`.
It is possible to have multiple connections for one type of database, just add this group again with new parameters.
```toml
delete_after = "7d"
delete_interval = "1h"
```


### delete_interval
Cleaning data of node, which are older than 7d.

```toml
delete_after = "7d"
```


### delete_interval
How often run the cleaning.

```toml
delete_interval = "1h"
```



## [[database.connection.influxdb]]
Save collected data to InfluxDB there would be the following measurements:
- node: store node spezific data i.e. clients memory, airtime
- global: store global data, i.e. count of clients and nodes
- firmware: store count of nodes tagged with firmware
- model: store count of nodes tagged with hardware model

```toml
enable   = false
address  = "http://localhost:8086"
database = "ffhb"
username = ""
password = ""
[database.connection.influxdb.tags]
site = "ffhb01"
system = "testing"
```


### enable
Enable the database connection instance to save collected values in a InfluxDB.

```toml
enable   = false
```


### address
Address to connect on InfluxDB server.

```toml
address  = "http://localhost:8086"
```


### database
Database on which the measurement should be stored.

```toml
database = "ffhb"
```


### username
Username to authenticate on InfluxDB

```toml
username = ""
```


### password
Password to authenticate on InfluxDB.

```toml
password = ""
```

### [database.connection.influxdb.tags]

You could set manuelle tags with inserting into a influxdb.
Usefull if you want to identify the yanic instance when you use multiple own on the same influxdb (e.g. multisites).

Warning, you could not overright tags which ware used by yanic (e.g. `nodeid`).
```toml
tagname = "value"
```




## [[database.connection.graphite]]
Save collected data to a graphite database.

```toml
enable   = false
address  = "localhost:2003"
prefix   = "freifunk"
```


### enable
Enable the database connection instance to save collected values in a graphite database.

```toml
enable   = false
```


### address
Address to connect on graphite server.

```toml
address = "localhost:2003"
```


### prefix
Prefix for every measurment key in this graphite database.

```toml
prefix = "freifunk"
```




## [[database.connection.logging]]
This database type is just for, debugging without a real database connection.
A example for other developers for new database types.

```toml
enable = false
path     = "/var/log/yanic.log"
```


### enable
Enable the database type logging.

```toml
enable = false
```


### path
Path to file where to store some examples with every line.

```toml
path     = "/var/log/yanic.log"
```
