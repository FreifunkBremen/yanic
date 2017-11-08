# Configuration

Here you would find a long description, maybe the description in [example file](https://github.com/FreifunkBremen/yanic/blob/master/config_example.toml) are enough for you.



{% method %}
## [respondd]
Group for configuration of respondd request.
{% sample lang="toml" %}
```toml
[respondd]
enable           = true
# synchronize    = "1m"
collect_interval = "1m"
interfaces       = ["br-ffhb"]
#port            = 10001
```
{% endmethod %}


{% method %}
### enable
Enable request and collection of data per respondd requests
{% sample lang="toml" %}
```toml
enable = true
```
{% endmethod %}


{% method %}
### synchronize
Delay startup until a multiple of the period since zero time
{% sample lang="toml" %}
```toml
synchronize      = "1m"
```
{% endmethod %}


{% method %}
### collect_interval
How often send request per respondd.

It will send UDP packets with multicast group `ff02::2:1001` and port `1001`.
If a node does not answer after the half time, it will request with the last know address under the port `1001`.
{% sample lang="toml" %}
```toml
collect_interval = "1m"
```
{% endmethod %}


{% method %}
### interfaces
Interface that has an IP in your mesh network
{% sample lang="toml" %}
```toml
interfaces       = ["br-ffhb"]
```
{% endmethod %}


{% method %}
### port
Define a port to listen and send the respondd packages.
If not set or set to 0 the kernel will use a random free port at its own.
{% sample lang="toml" %}
```toml
port              = 10001
```
{% endmethod %}



{% method %}
## [webserver]
Yanic has a little build-in webserver, which statically serves a directory.
This is useful for testing purposes or for a little standalone installation.
{% sample lang="toml" %}
```toml
[webserver]
enable  = false
bind    = "127.0.0.1:8080"
webroot = "/var/www/html/meshviewer"
```
{% endmethod %}


{% method %}
### enable
Enable to start the built-in webserver of Yanic
{% sample lang="toml" %}
```toml
enable  = false
```
{% endmethod %}


{% method %}
### bind
On which ip address and port listen the webserver
{% sample lang="toml" %}
```toml
bind    = "127.0.0.1:8080"
```
{% endmethod %}


{% method %}
### webroot
The path to a folder, which files are published on this webserver.
{% sample lang="toml" %}
```toml
webroot = "/var/www/html/meshviewer"
```
{% endmethod %}



{% method %}
## [nodes]
{% sample lang="toml" %}
```toml
[nodes]
enable         = true
state_path     = "/var/lib/yanic/state.json"
save_interval = "5s"
offline_after = "10m"
prune_after = "7d"
```
{% endmethod %}


{% method %}
### enable
Enable the storing and writing of json files.
{% sample lang="toml" %}
```toml
enable         = true
```
{% endmethod %}


{% method %}
### state_path
A json file to cache all data collected directly from respondd.
{% sample lang="toml" %}
```toml
state_path     = "/var/lib/yanic/state.json"
```
{% endmethod %}


{% method %}
### save_interval
Export nodes and graph periodically.
{% sample lang="toml" %}
```toml
save_interval = "5s"
```
{% endmethod %}


{% method %}
### offline_after
Set node to offline if not seen within this period.
{% sample lang="toml" %}
```toml
offline_after = "10m"
```
{% endmethod %}


{% method %}
### prune_after
Prune data in RAM, cache-file and output json files (i.e. nodes.json) that were inactive for longer than.
{% sample lang="toml" %}
```toml
prune_after = "7d"
```
{% endmethod %}



{% method %}
## [meshviewer]
{% sample lang="toml" %}
```toml
[meshviewer]
version        = 2
nodes_path     = "/var/www/html/meshviewer/data/nodes.json"
graph_path     = "/var/www/html/meshviewer/data/graph.json"
```
{% endmethod %}


{% method %}
### version
The structure version of the output which should be generated (i.e. nodes.json)
* version 1 is accepted by the legacy meshviewer (which is the master branch)
* https://github.com/ffnord/meshviewer/tree/master
* version 2 is accepted by the new version of meshviewer (which are in legacy develop branch or newer)
* https://github.com/ffnord/meshviewer/tree/dev
* https://github.com/ffrgb/meshviewer/tree/develop
{% sample lang="toml" %}
```toml
version = 2
```
{% endmethod %}


{% method %}
### nodes_path
The path, where to store nodes.json (supports version 1 and two, see `nodes_version`)
{% sample lang="toml" %}
```toml
nodes_path = "/var/www/html/meshviewer/data/nodes.json"
```
{% endmethod %}


{% method %}
### graph_path
The path, where to store graph.json (only version 1)
{% sample lang="toml" %}
```toml
graph_path = "/var/www/html/meshviewer/data/graph.json"
```
{% endmethod %}



{% method %}
## [database]
The database organize all database types.
For all database types the is a internal job, which reset data for nodes (global statistics are still stored).
_(We have for privacy policy to store node data for maximum seven days.)_

Every database type has his own configuration under `database.connection`.
It is possible to have multiple connections for one type of database, just add this group again with new parameters.
{% sample lang="toml" %}
```toml
delete_after = "7d"
delete_interval = "1h"
```
{% endmethod %}


{% method %}
### delete_interval
This will send delete commands to the database to prune data which is older than:
{% sample lang="toml" %}
```toml
delete_after = "7d"
```
{% endmethod %}


{% method %}
### delete_interval
How often run the delete commands.
{% sample lang="toml" %}
```toml
delete_interval = "1h"
```
{% endmethod %}


{% method %}
## [[database.connection.influxdb]]
Save collected data to InfluxDB.
There are would be the following measurements:
- node: store node specific data i.e. clients memory, airtime
- global: store global data, i.e. count of clients and nodes
- firmware: store the count of nodes tagged with firmware
- model: store the count of nodes tagged with hardware model
{% sample lang="toml" %}
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
{% endmethod %}


{% method %}
### enable
Enable the database connection instance to save collected values in a InfluxDB.
{% sample lang="toml" %}
```toml
enable   = false
```
{% endmethod %}


{% method %}
### address
Address to connect on InfluxDB server.
{% sample lang="toml" %}
```toml
address  = "http://localhost:8086"
```
{% endmethod %}


{% method %}
### database
Database on which the measurement should be stored.
{% sample lang="toml" %}
```toml
database = "ffhb"
```
{% endmethod %}


{% method %}
### username
Username to authenticate on InfluxDB
{% sample lang="toml" %}
```toml
username = ""
```
{% endmethod %}


{% method %}
### password
Password to authenticate on InfluxDB.
{% sample lang="toml" %}
```toml
password = ""
```
{% endmethod %}


{% method %}
### [database.connection.influxdb.tags]

You could set manuelle tags with inserting into a influxdb.
Usefull if you want to identify the yanic instance when you use multiple own on the same influxdb (e.g. multisites).

Warning:
Tags used by Yanic would override the tags from this config (e.g. `nodeid`, `hostname`, `owner`, `model`, `firmware_base`, `firmware_release`, `frequency11g`, `frequency11a`).
{% sample lang="toml" %}
```toml
tagname1 = "tagvalue 1s"
# some usefull e.g.:
system   = "productive"
site     = "ffhb"
```
{% endmethod %}



{% method %}
## [[database.connection.graphite]]
Save collected data to a graphite database.
{% sample lang="toml" %}
```toml
enable   = false
address  = "localhost:2003"
prefix   = "freifunk"
```
{% endmethod %}


{% method %}
### enable
Enable the database connection instance to save collected values in a graphite database.
{% sample lang="toml" %}
```toml
enable   = false
```
{% endmethod %}


{% method %}
### address
Address to connect on graphite server.
{% sample lang="toml" %}
```toml
address = "localhost:2003"
```
{% endmethod %}


{% method %}
### prefix
Graphite is replacing every "." in the metric name with a slash "/" and uses
that for the file system hierarchy it generates. it is recommended to at least
move the metrics out of the root namespace (that would be the empty prefix).
If you only intend to run one community and only freifunk on your graphite node
then the prefix can be set to anything (including the empty string) since you
probably wont care much about "polluting" the namespace.
{% sample lang="toml" %}
```toml
prefix = "freifunk"
```
{% endmethod %}



{% method %}
## [[database.connection.logging]]
This database type is just for, debugging without a real database connection.
A example for other developers for new database types.
{% sample lang="toml" %}
```toml
enable = false
path     = "/var/log/yanic.log"
```
{% endmethod %}


{% method %}
### enable
Enable the database type logging.
{% sample lang="toml" %}
```toml
enable = false
```
{% endmethod %}


{% method %}
### path
Path to file where to store some examples with every line.
{% sample lang="toml" %}
```toml
path     = "/var/log/yanic.log"
```
{% endmethod %}
