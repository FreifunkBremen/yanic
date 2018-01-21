# Configuration

Here you would find a long description, maybe the description in [example file](https://github.com/FreifunkBremen/yanic/blob/master/config_example.toml) are enough for you.

The config file for Yanic written in "Tom's Obvious, Minimal Language." [syntax](https://github.com/toml-lang/toml).
(if you need somethink multiple times, checkout out the [[array of table]] section)

## [respondd]
{% method %}
Group for configuration of respondd request.
{% sample lang="toml" %}
```toml
[respondd]
enable           = true
# synchronize    = "1m"
collect_interval = "1m"
interfaces       = ["br-ffhb"]
sites            = ["ffhb"]
#port            = 10001
```
{% endmethod %}


### enable
{% method %}
Enable request and collection of data per respondd requests
{% sample lang="toml" %}
```toml
enable = true
```
{% endmethod %}


### synchronize
{% method %}
Delay startup until a multiple of the period since zero time
{% sample lang="toml" %}
```toml
synchronize      = "1m"
```
{% endmethod %}


### collect_interval
{% method %}
How often send request per respondd.

It will send UDP packets with multicast group `ff02::2:1001` and port `1001`.
If a node does not answer after the half time, it will request with the last know address under the port `1001`.
{% sample lang="toml" %}
```toml
collect_interval = "1m"
```
{% endmethod %}


### interfaces
{% method %}
Interface that has an IP in your mesh network
{% sample lang="toml" %}
```toml
interfaces       = ["br-ffhb"]
```
{% endmethod %}


### sites
{% method %}
List of sites to save stats for (empty for global only)
{% sample lang="toml" %}
```toml
sites            = ["ffhb"]
```
{% endmethod %}


### port
{% method %}
Define a port to listen and send the respondd packages.
If not set or set to 0 the kernel will use a random free port at its own.
{% sample lang="toml" %}
```toml
port              = 10001
```
{% endmethod %}



## [webserver]
{% method %}
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


### enable
{% method %}
Enable to start the built-in webserver of Yanic
{% sample lang="toml" %}
```toml
enable  = false
```
{% endmethod %}


### bind
{% method %}
On which ip address and port listen the webserver
{% sample lang="toml" %}
```toml
bind    = "127.0.0.1:8080"
```
{% endmethod %}


### webroot
{% method %}
The path to a folder, which files are published on this webserver.
{% sample lang="toml" %}
```toml
webroot = "/var/www/html/meshviewer"
```
{% endmethod %}



## [nodes]
{% method %}
{% sample lang="toml" %}
```toml
[nodes]
state_path     = "/var/lib/yanic/state.json"
prune_after    = "7d"
save_interval  = "5s"
offline_after  = "10m"
```
{% endmethod %}


### state_path
{% method %}
A json file to cache all data collected directly from respondd.
{% sample lang="toml" %}
```toml
state_path     = "/var/lib/yanic/state.json"
```
{% endmethod %}


### prune_after
{% method %}
Prune data in RAM, cache-file and output json files (i.e. nodes.json) that were inactive for longer than.
{% sample lang="toml" %}
```toml
prune_after = "7d"
```
{% endmethod %}


### save_interval
{% method %}
Export nodes and graph periodically.
{% sample lang="toml" %}
```toml
save_interval = "5s"
```
{% endmethod %}


### offline_after
{% method %}
Set node to offline if not seen within this period.
{% sample lang="toml" %}
```toml
offline_after = "10m"
```
{% endmethod %}


## [[nodes.output.example]]
{% method %}
This example block shows all option which is useable for every following output type.
Every output type has his own configuration under `nodes.output`.
It is possible to have multiple output for one type of output, just add this group again with new parameters (see toml [[array of table]]).
{% sample lang="toml" %}
```toml
[[nodes.output.example]]
enable = true
[nodes.output.example.filter]
no_owner  = true
blacklist = ["00112233445566", "1337f0badead"]
sites = ["ffhb"]
has_location = true
[nodes.output.example.filter.in_area]
latitude_min  = 34.30
latitude_max  = 71.85
longitude_min = -24.96
longitude_max = 39.72
```
{% endmethod %}

### enable
{% method %}
Each output format has its own config block and needs to be enabled by adding:
{% sample lang="toml" %}
```toml
enable = true
```
{% endmethod %}

### [nodes.output.example.filter]
{% method %}
For each output format there can be set different filters
{% sample lang="toml" %}
```toml
[nodes.output.example.filter]
no_owner  = true
blacklist = ["00112233445566", "1337f0badead"]
sites = ["ffhb"]
has_location = true
[nodes.output.example.filter.in_area]
latitude_min  = 34.30
latitude_max  = 71.85
longitude_min = -24.96
longitude_max = 39.72
```
{% endmethod %}


### no_owner
{% method %}
Set to false, if you want the json files to contain the owner information


**WARNING: if it is not set, it will publish contact information of other persons.**

{% sample lang="toml" %}
```toml
no_owner = true
```
{% endmethod %}


### blacklist
{% method %}
List of nodeids of nodes that should be filtered out, so they won't appear in output
{% sample lang="toml" %}
```toml
blacklist = ["00112233445566", "1337f0badead"]
```
{% endmethod %}


### sites
{% method %}
List of site_codes of nodes that should be included in output
{% sample lang="toml" %}
```toml
sites = ["ffhb"]
```
{% endmethod %}


### has_location
{% method %}
set has_location to true if you want to include only nodes that have geo-coordinates set
(setting this to false has no sensible effect, unless you'd want to hide nodes that have coordinates)
{% sample lang="toml" %}
```toml
has_location = true
```
{% endmethod %}


### [nodes.output.example.filter.in_area]
{% method %}
nodes outside this area are not shown on the map but are still listed as a node without coordinates
{% sample lang="toml" %}
```toml
latitude_min = 34.30
latitude_max = 71.85
longitude_min = -24.96
longitude_max = 39.72
```
{% endmethod %}



## [[nodes.output.meshviewer-ffrgb]]
{% method %}
The new json file format for the [meshviewer](https://github.com/ffrgb/meshviewer) developed in Regensburg.

{% sample lang="toml" %}
```toml
[[nodes.output.meshviewer-ffrgb]]
enable   = true
path     = "/var/www/html/meshviewer/data/meshviewer.json"
#[nodes.output.meshviewer-ffrgb.filter]
#no_owner = false
#blacklist = ["00112233445566", "1337f0badead"]
#has_location = true

#[nodes.output.meshviewer-ffrgb.filter.in_area]
#latitude_min = 34.30
#latitude_max = 71.85
#longitude_min = -24.96
#longitude_max = 39.72
```
{% endmethod %}


### path
{% method %}
The path, where to store meshviewer.json
{% sample lang="toml" %}
```toml
path     = "/var/www/html/meshviewer/data/meshviewer.json"
```
{% endmethod %}



## [[nodes.output.meshviewer]]
{% method %}
{% sample lang="toml" %}
```toml
[[nodes.output.meshviewer]]
enable         = false
version        = 2
nodes_path     = "/var/www/html/meshviewer/data/nodes.json"
graph_path     = "/var/www/html/meshviewer/data/graph.json"
```
{% endmethod %}


### version
{% method %}
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


### nodes_path
{% method %}
The path, where to store nodes.json (supports version 1 and two, see `nodes_version`)
{% sample lang="toml" %}
```toml
nodes_path = "/var/www/html/meshviewer/data/nodes.json"
```
{% endmethod %}


### graph_path
{% method %}
The path, where to store graph.json (only version 1)
{% sample lang="toml" %}
```toml
graph_path = "/var/www/html/meshviewer/data/graph.json"
```
{% endmethod %}



## [[nodes.output.nodelist]]
{% method %}
The nodelist output is a minimal output with current state of collected data.
Should be preferred to use it on the [ffapi](https://freifunk.net/api-generator/) for the [freifunk-karte.de](https://freifunk-karte.de)
{% sample lang="toml" %}
```toml
[[nodes.output.nodelist]]
enable   = false
path     = "/var/www/html/meshviewer/data/nodelist.json"
#[nodes.output.nodelist.filter]
#no_owner = false
```
{% endmethod %}


### path
{% method %}
The path, where to store nodelist.json
{% sample lang="toml" %}
```toml
path     = "/var/www/html/meshviewer/data/nodelist.json"
```
{% endmethod %}



## [database]
{% method %}
The database organize all database types.
For all database types the is a internal job, which reset data for nodes (global statistics are still stored).
_(We have for privacy policy to store node data for maximum seven days.)_
{% sample lang="toml" %}
```toml
delete_after = "7d"
delete_interval = "1h"
```
{% endmethod %}


### delete_after
{% method %}
This will send delete commands to the database to prune data which is older than:
{% sample lang="toml" %}
```toml
delete_after = "7d"
```
{% endmethod %}


### delete_interval
{% method %}
How often run the delete commands.
{% sample lang="toml" %}
```toml
delete_interval = "1h"
```
{% endmethod %}


## [[database.connection.example]]
{% method %}
This example block shows all option which is useable for every following database type.
Every database type has his own configuration under `database.connection`.
It is possible to have multiple connections for one type of database, just add this group again with new parameters (see toml [[array of table]]).
{% sample lang="toml" %}
```toml
[[database.connection.example]]
enable = true
```
{% endmethod %}


### enable
{% method %}
Each database-connection has its own config block and needs to be enabled by adding:
{% sample lang="toml" %}
```toml
enable = true
```
{% endmethod %}



## [[database.connection.influxdb]]
{% method %}
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
tagname1 = "tagvalue 1"
system   = "productive"
site     = "ffhb"
```
{% endmethod %}


### address
{% method %}
Address to connect on InfluxDB server.
{% sample lang="toml" %}
```toml
address  = "http://localhost:8086"
```
{% endmethod %}


### database
{% method %}
Database on which the measurement should be stored.
{% sample lang="toml" %}
```toml
database = "ffhb"
```
{% endmethod %}


### username
{% method %}
Username to authenticate on InfluxDB
{% sample lang="toml" %}
```toml
username = ""
```
{% endmethod %}


### password
{% method %}
Password to authenticate on InfluxDB.
{% sample lang="toml" %}
```toml
password = ""
```
{% endmethod %}


### [database.connection.influxdb.tags]
{% method %}
You could set manuelle tags with inserting into a influxdb.
Useful if you want to identify the yanic instance when you use multiple own on the same influxdb (e.g. multisites).

Warning:
Tags used by Yanic would override the tags from this config (e.g. `nodeid`, `hostname`, `owner`, `model`, `firmware_base`, `firmware_release`, `frequency11g`, `frequency11a`).
{% sample lang="toml" %}
```toml
tagname1 = "tagvalue 1s"
# some useful e.g.:
system   = "productive"
site     = "ffhb"
```
{% endmethod %}



## [[database.connection.graphite]]
{% method %}
Save collected data to a graphite database.
{% sample lang="toml" %}
```toml
enable   = false
address  = "localhost:2003"
prefix   = "freifunk"
```
{% endmethod %}


### address
{% method %}
Address to connect on graphite server.
{% sample lang="toml" %}
```toml
address = "localhost:2003"
```
{% endmethod %}


### prefix
{% method %}
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



## [[database.connection.logging]]
{% method %}
This database type is just for, debugging without a real database connection.
A example for other developers for new database types.
{% sample lang="toml" %}
```toml
enable = false
path     = "/var/log/yanic.log"
```
{% endmethod %}


### path
{% method %}
Path to file where to store some examples with every line.
{% sample lang="toml" %}
```toml
path     = "/var/log/yanic.log"
```
{% endmethod %}
