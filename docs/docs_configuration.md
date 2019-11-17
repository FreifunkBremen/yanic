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

#[respondd.sites.example]
#domains            = ["city"]

#[[respondd.custom_field]]
#name = zip
# You can use arbitrary GJSON expressions here, see https://github.com/tidwall/gjson
# We expect this expression to return a string.
#path = nodeinfo.location.zip

[[respondd.interfaces]]
ifname             = "br-ffhb"
#ip_address        = "fe80::..."
#send_no_request   = false
#multicast_address = "ff02::2:1001"
#port              = 10001
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

It will send UDP packets with multicast address `ff05::2:1001` and port `1001`.
If a node does not answer after the half time, it will request with the last know address under the port `1001`.
{% sample lang="toml" %}
```toml
collect_interval = "1m"
```
{% endmethod %}


### [respondd.sites.example]
{% method %}
Tables of sites to save stats for (not exists for global only).
Here is the site _ffhb_.
{% sample lang="toml" %}
```toml
[respondd.sites.ffhb]
domains            = ["city"]
```
{% endmethod %}

#### domains
{% method %}
list of domains on this site to save stats for (empty for global only)
{% sample lang="toml" %}
```toml
domains            = ["city"]
```
{% endmethod %}


### [[respondd.interfaces]]
{% method %}
Interface that has an ip address in your mesh network.
It is possible to have multiple interfaces, just add this group again with new parameters (see toml [[array of table]]).
{% sample lang="toml" %}
```toml
[[respondd.interfaces]]
ifname             = "br-ffhb"
#ip_address        = "fe80::..."
#send_no_request   = false
#multicast_address = "ff02::2:1001"
#port              = 10001
```
{% endmethod %}

### ifname
{% method %}
name of interface on which this collector is running.
{% sample lang="toml" %}
```toml
ifname              = "br-ffhb"
```
{% endmethod %}

### ip_address
{% method %}
ip address is the own address which is used for sending.
If not set or set with empty string it will take an address of ifname.
(It prefers the link local address, so at babel mesh-network it should be configurated)
{% sample lang="toml" %}
```toml
ip_address          = "fe80::..."
```
{% endmethod %}

### send_no_request
{% method %}
Disable sending multicast respondd request.
For receiving only respondd packages e.g. database respondd.
{% sample lang="toml" %}
```toml
send_no_request     = true
```
{% endmethod %}

### multicast_address
{% method %}
Multicast address to destination of respondd.
If not set or set with empty string it will take the batman default multicast address `ff05::2:1001`
(Needed to set for legacy `ff02::2:1001`)
{% sample lang="toml" %}
```toml
multicast_address    = "ff02::2:1001"
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

### [[respondd.custom_fields]]
{% method %}
If you have custom respondd fields, you can ask Yanic to also collect these.

NOTE: This does not automatically include these fields in the output.
The meshviewer-ffrgb output module will include them under "custom_fields",
but other modules may simply ignore them.

{% sample lang="toml" %}
```toml
name = zip
# You can use arbitrary GJSON expressions here, see https://github.com/tidwall/gjson
# We expect this expression to return a string.
path = nodeinfo.location.zip
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
domain_as_site = true
domain_append_site = true
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

### domain_as_site
{% method %}
Replace the `site_code` with the `domain_code` in this output.
e.g. `site_code='ffhb',domain_code='city'` becomes `site_code='city', domain_code=''`
{% sample lang="toml" %}
```toml
domain_as_site = true
```
{% endmethod %}

### domain_append_site
{% method %}
Append on the `site_code` the `domain_code` with a `.` in this output.
e.g. `site_code='ffhb',domain_code='city'` becomes `site_code='ffhb.city', domain_code=''`
{% sample lang="toml" %}
```toml
domain_append_site = true
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



## [[nodes.output.geojson]]
{% method %}
The geojson output produces a geojson file which contains the location data of all monitored nodes to be used to visualize the location of the nodes.
It is optimized to be used with [UMap](https://github.com/umap-project/umap) but should work with other tools as well.

Here is a public demo provided by Freifunk Muenchen: http://u.osmfr.org/m/328494/
{% sample lang="toml" %}
```toml
[[nodes.output.geojson]]
enable   = true
path = "/var/www/html/meshviewer/data/nodes.geojson"
```
{% endmethod %}


### path
{% method %}
The path, where to store nodes.geojson
{% sample lang="toml" %}
```toml
path     = "/var/www/html/meshviewer/data/nodes.geojson"
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

## [[nodes.output.raw]]
{% method %}
This output takes the respondd response as sent by the node and includes it in a JSON document.
{% endmethod %}


### path
{% method %}
The path, where to store raw.json
{% sample lang="toml" %}
```toml
path     = "/var/www/html/meshviewer/data/raw.json"
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
insecure_skip_verify = false
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

### insecure_skip_verify
{% method %}
Skip insecure verify for self-signed certificates.
{% sample lang="toml" %}
```toml
insecure_skip_verify = true
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



## [[database.connection.respondd]]
{% method %}
Forward collected respondd package to a address
(e.g. to another respondd collector like a central yanic instance or hopglass)
{% sample lang="toml" %}
```toml
enable   = false
type     = "udp6"
address  = "stats.bremen.freifunk.net:11001"
```
{% endmethod %}


### type
{% method %}
Type of network to create a connection.

Known networks are "tcp", "tcp4" (IPv4-only), "tcp6" (IPv6-only), "udp", "udp4" (IPv4-only), "udp6" (IPv6-only), "ip", "ip4" (IPv4-only), "ip6" (IPv6-only), "unix", "unixgram" and "unixpacket".
{% sample lang="toml" %}
```toml
type     = "udp6"
```
{% endmethod %}


### address
{% method %}
Destination address to connect/send respondd package.
{% sample lang="toml" %}
```toml
address  = "stats.bremen.freifunk.net:11001"
```
{% endmethod %}



## [[database.connection.logging]]
{% method %}
This database type is just for, debugging without a real database connection.
A example for other developers for new database types.
{% sample lang="toml" %}
```toml
enable   = false
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
