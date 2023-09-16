# Configuration

Here you would find a long description, maybe the description in [example file](https://github.com/FreifunkBremen/yanic/blob/main/config_example.toml) are enough for you.

The config file for Yanic written in "Tom's Obvious, Minimal Language." [syntax](https://github.com/toml-lang/toml).
(if you need somethink multiple times, checkout out the [[array of table]] section)

## [respondd]
Group for configuration of respondd request.

```toml
# Send respondd request to update information
[respondd]
enable           = true # (1)
# Delay startup until a multiple of the period since zero time
synchronize    = "1m" # (2)
# how often request per multicast
collect_interval = "1m" # (3)

# If you have custom respondd fields, you can ask Yanic to also collect these.
#[[respondd.custom_field]] (4)
#name = zip
# You can use arbitrary GJSON expressions here, see https://github.com/tidwall/gjson
# We expect this expression to return a string.
#path = nodeinfo.location.zip

# table of a site to save stats for (not exists for global only)
#[respondd.sites.example] (5)
## list of domains on this site to save stats for (empty for global only)
#domains            = ["city"] (6)


# interface that has an IP in your mesh network
[[respondd.interfaces]] # (7)
# name of interface on which this collector is running
ifname             = "br-ffhb" # (8) 
# ip address which is used for sending
# (optional - without definition used a address of ifname - preferred link local)
ip_address        = "fd2f:5119:f2d::5" # (9)
# disable sending multicast respondd request
# (for receiving only respondd packages e.g. database respondd)
#send_no_request   = false (10)
# multicast address to destination of respondd
# (optional - without definition used default ff05::2:1001)
# Very old gluon uses "ff02::2:1001" as multicast, newer use ff05::2:1001. If you have old and new
# gluon nodes on the same network, create a separate "respondd.interfaces" section for each mutlicast address.
#multicast_address = "ff02::2:1001" (11)
# define a port to listen
# if not set or set to 0 the kernel will use a random free port at its own
#port              = 10001 (12)
```

1.  Enable request and collection of data per respondd requests

2.  Delay startup until a multiple of the period since zero time

3.  How often send request per respondd.

    It will send UDP packets with multicast address `ff05::2:1001` and port `1001`.
    If a node does not answer after the half time, it will request with the last know address under the port `1001`.

4.  If you have custom respondd fields, you can ask Yanic to also collect these.

    It is possible to have multiple custom fields, just add this group again with new parameters (see toml [[array of table]]).

    !!! info
        This does not automatically include these fields in the output.
        The meshviewer-ffrgb output module will include them under "custom_fields",
        but other modules may simply ignore them.


5.  Tables of sites to save stats for (not exists for global only).
    Here is the site _ffhb_.

    !!! example
        ```toml
        [respondd.sites.ffhb]
        domains = ["city"]
        ```

6.  list of domains on this site to save stats for (empty for global only)

7.  Interface that has an ip address in your mesh network.

    It is possible to have multiple interfaces, just add this group again with new parameters (see toml [[array of table]]).

8.  name of interface on which this collector is running.

9.  ip address is the own address which is used for sending.
    If not set or set with empty string it will take an address of ifname.
    (It prefers the link local address, so at babel mesh-network it should be configurated)

10. Disable sending multicast respondd request.
    For receiving only respondd packages e.g. database respondd.

11. Multicast address to destination of respondd.
    If not set or set with empty string it will take the batman default multicast address `ff05::2:1001`
    (Needed to set for legacy `ff02::2:1001`)

12. Define a port to listen and send the respondd packages.
    If not set or set to 0 the kernel will use a random free port at its own.

## [webserver]
Yanic has a little build-in webserver, which statically serves a directory.
This is useful for testing purposes or for a little standalone installation.

```toml
# A little build-in webserver, which statically serves a directory.
# This is useful for testing purposes or for a little standalone installation.
[webserver]
enable  = false # (1)
bind    = "127.0.0.1:8080" # (2)
webroot = "/var/www/html/meshviewer" # (3)
```

1.  Enable to start the built-in webserver of Yanic
2.  On which ip address and port listen the webserver
3.  The path to a folder, which files are published on this webserver.

## [nodes]
```toml
[nodes]
# Cache file
# a json file to cache all data collected directly from respondd
state_path     = "/var/lib/yanic/state.json" # (1)
# prune data in RAM, cache-file and output json files (i.e. nodes.json)
# that were inactive for longer than
prune_after    = "7d" # (2)
# Export nodes and graph periodically
save_interval  = "5s" # (3)
# Set node to offline if not seen within this period
offline_after  = "10m" # (4)
```

1.  A json file to cache all data collected directly from respondd.
2.  Prune data in RAM, cache-file and output json files (i.e. nodes.json) that were inactive for longer than.
3.  Export nodes and graph periodically.
4.  Set node to offline if not seen within this period.


## [[nodes.output.example]]
This example block shows all option which is useable for every following output type.
Every output type has his own configuration under `nodes.output`.
It is possible to have multiple output for one type of output, just add this group again with new parameters (see toml [[array of table]]).

```toml
[[nodes.output.example]]
enable = true # (1)

[nodes.output.example.filter] # (2)
no_owner  = true # if it is not set, it will publish contact information of other persons (3)
blocklist = ["00112233445566", "1337f0badead"] # (4)
sites = ["ffhb"] # (5)
domain_as_site = true # (6)
domain_append_site = true # (7)
has_location = true # (8)

[nodes.output.example.filter.in_area] # (9)
latitude_min  = 34.30
latitude_max  = 71.85
longitude_min = -24.96
longitude_max = 39.72
```

1.  Each output format has its own config block and needs to be enabled by adding:
2.  For each output format there can be set different filters
3.  Set to false, if you want the json files to contain the owner information

    !!! warning
        if it is not set, it will publish contact information of other persons.

4.  List of nodeids of nodes that should be filtered out, so they won't appear in output

5.  List of site_codes of nodes that should be included in output

6.  Replace the `site_code` with the `domain_code` in this output.
    
    e.g. `site_code='ffhb',domain_code='city'` becomes `site_code='city', domain_code=''`

7.  Append on the `site_code` the `domain_code` with a `.` in this output.

    e.g. `site_code='ffhb',domain_code='city'` becomes `site_code='ffhb.city', domain_code=''`

8.  set has_location to true if you want to include only nodes that have geo-coordinates set

    (setting this to false has no sensible effect, unless you'd want to hide nodes that have coordinates)

9.  nodes outside this area are not shown on the map but are still listed as a node without coordinates



### [[nodes.output.geojson]]
The geojson output produces a geojson file which contains the location data of all monitored nodes to be used to visualize the location of the nodes.
It is optimized to be used with [UMap](https://github.com/umap-project/umap) but should work with other tools as well.

Here is a public demo provided by Freifunk Muenchen: http://u.osmfr.org/m/328494/

```toml hl_lines="1-3"
[[nodes.output.geojson]]
enable   = true
path = "/var/www/html/meshviewer/data/nodes.geojson" # (1)
[nodes.output.geojson.filter]
no_owner = false
```

1. The path, where to store nodes.geojson



### [[nodes.output.meshviewer-ffrgb]]
The new json file format for the [meshviewer](https://github.com/ffrgb/meshviewer) developed in Regensburg.

```toml hl_lines="1-3"
[[nodes.output.meshviewer-ffrgb]]
enable   = true
path     = "/var/www/html/meshviewer/data/meshviewer.json" # (1)
# like on every output, here some filters, for example using this block:
[nodes.output.meshviewer-ffrgb.filter]
no_owner = false
blocklist = ["00112233445566", "1337f0badead"]

[nodes.output.meshviewer-ffrgb.filter.in_area]
latitude_min = 34.30
latitude_max = 71.85
longitude_min = -24.96
longitude_max = 39.72
```

1. The path, where to store meshviewer.json



### [[nodes.output.meshviewer]]
```toml hl_lines="1-5"
[[nodes.output.meshviewer]]
enable         = false
version        = 2 # (1)
nodes_path     = "/var/www/html/meshviewer/data/nodes.json" # (3)
graph_path     = "/var/www/html/meshviewer/data/graph.json" # (4)
[nodes.output.meshviewer.filter]
no_owner = false
```

1.  The structure version of the output which should be generated (i.e. nodes.json)

    * version `1` is accepted by the legacy meshviewer (which is the master branch)
        * https://github.com/ffnord/meshviewer/tree/master
     <!-- bug: count as 2 ... -->
    * version `2` is accepted by the new version of meshviewer (which are in legacy develop branch or newer)
        * https://github.com/ffnord/meshviewer/tree/dev
        * https://github.com/ffrgb/meshviewer/tree/develop
3.  The path, where to store nodes.json (supports version 1 and two, see `nodes_version`)
4.  The path, where to store graph.json (only version 1)



### [[nodes.output.nodelist]]
The nodelist output is a minimal output with current state of collected data.
Should be preferred to use it on the [ffapi](https://freifunk.net/api-generator/) for the [freifunk-karte.de](https://freifunk-karte.de)

```toml hl_lines="1-3"
[[nodes.output.nodelist]]
enable   = false
path     = "/var/www/html/meshviewer/data/nodelist.json" # (1)
[nodes.output.nodelist.filter]
no_owner = false
```

1.  The path, where to store nodelist.json



### [[nodes.output.prometheus-sd]]
The Prometheus Service Discovery (SD) output is a output with the list of addresses of the nodes to use them in later exporter by prometheus.
For usage in Prometheus read there Documentation [Use file-based service discovery to discover scrape targets](https://prometheus.io/docs/guides/file-sd/).

```toml
[[nodes.output.prometheus-sd]]
enable         = false
path           = "/var/www/html/meshviewer/data/prometheus-sd.json" # (1)
# ip = lates recieved ip, node_id = node id from host
target_address = "ip" # (2)

# Labels of the data (optional)
[nodes.output.prometheus-sd.labels] # (3)
labelname1 = "labelvalue 1"
# some useful e.g.:
hosts   = "ffhb"
service = "yanic"
```

1.  The path, where to store prometheus-sd.json
2.  In the prometheus-sd.json the usage of which information of the node as targets (address).
    
    Use the `node_id` as value, to put the Node ID into the target list as address.
    
    Use the `ip` as value to put the last IP address into the target list from where the respondd message is recieved (maybe a link-local address).
    
    Default value is `ip`.

3.  You could optional set manuelle labels with inserting into a prometheus-sd.json.
    Useful if you want to identify the yanic instance when you use multiple own on the same prometheus database (e.g. multisites).

### [[nodes.output.raw]]
This output takes the respondd response as sent by the node and includes it in a JSON document.
```toml hl_lines="1-3"
[[nodes.output.raw]]
enable   = false
path     = "/var/www/html/meshviewer/data/raw.json" # (1)
[nodes.output.raw.filter]
no_owner = false
```

1.  The path, where to store raw.json


### [[nodes.output.raw-jsonl]]
This output takes the respondd response as sent by the node and inserts it into a line-separated JSON document (JSONL). In this format, each line can be interpreted as a separate JSON element, which is useful for json streaming. The first line is a json object containing the timestamp and version of the file. This is followed by a line for each node, each containing a json object.
```toml hl_lines="1-3"
[[nodes.output.raw-jsonl]]
enable   = false
path     = "/var/www/html/meshviewer/data/raw.jsonl" # (1)
[nodes.output.raw-jsonl.filter]
no_owner = false
```

1.  The path, where to store raw.jsonl

## [database]
The database organize all database types.
For all database types the is a internal job, which reset data for nodes (global statistics are still stored).
_(We have for privacy policy to store node data for maximum seven days.)_

```toml
[database]
delete_after = "7d" # (1)
delete_interval = "1h" # (2)
```

1. This will send delete commands to the database to prune data which is older than:
2. How often run the delete commands.


## [[database.connection.example]]
This example block shows all option which is useable for every following database type.
Every database type has his own configuration under `database.connection`.
It is possible to have multiple connections for one type of database, just add this group again with new parameters (see toml [[array of table]]).

```toml
[[database.connection.example]]
enable = true # (1)
```

1. Each database-connection has its own config block and needs to be enabled by adding:

### [[database.connection.influxdb]]
Save collected data to InfluxDB.
There are would be the following measurements:
- node: store node specific data i.e. clients memory, airtime
- link: store link tq between two interfaces of two different nodes
- global: store global data, i.e. count of clients and nodes
- firmware: store the count of nodes tagged with firmware
- model: store the count of nodes tagged with hardware model
- autoupdater: store the count of autoupdate branch

```toml
[[database.connection.influxdb]]
enable   = false
address  = "http://localhost:8086" # (1)
database = "ffhb" # (2)
username = "" # (3)
password = "" # (4)
# insecure_skip_verify = true (5)

[database.connection.influxdb.tags] # (6)
tagname1 = "tagvalue 1"
system   = "productive"
site     = "ffhb"
```

1.  Address to connect on InfluxDB server.
2.  Database on which the measurement should be stored.
3.  Username to authenticate on InfluxDB
4.  Password to authenticate on InfluxDB.
5.  Skip insecure verify for self-signed certificates.
6.  You could set manuelle tags with inserting into a influxdb.

    Useful if you want to identify the yanic instance when you use multiple own on the same influxdb (e.g. multisites).

    !!! warning
        Tags used by Yanic would override the tags from this config (e.g. `nodeid`, `hostname`, `owner`, `model`, `firmware_base`, `firmware_release`, `frequency11g`, `frequency11a`).

### [[database.connection.influxdb2]]
Save collected data to InfluxDB2.

There are the following measurments:

  - **node**: store node specific data i.e. clients memory, airtime
  - **link**: store link tq between two interfaces of two different nodes with i.e. nodeid, address, hostname
  - **global**: store global data, i.e. count of clients and nodes
  - **firmware**: store the count of nodes tagged with firmware
  - **model**: store the count of nodes tagged with hardware model
  - **autoupdater**: store the count of autoupdate branch


!!! info
    A bucket has to be set in buckets and buchet_default otherwise yanic would panic.

!!! warning
    yanic do NOT prune node's data (so please setup it in InfluxDB2 setup).

    We highly recommend to setup e.g. [Data retention](https://docs.influxdata.com/influxdb/v2/reference/internals/data-retention/) in your InfluxDB2 server per measurements.

```toml
[[database.connection.influxdb2]]
enable   = false
address  = "http://localhost:8086" # (1)
token = "" # (2)
organization_id = "" # (3)
bucket_default = "" # (4)

[database.connection.influxdb2.buckets] # (5)
#link = "yanic-temp"
#node = "yanic-temp"
#dhcp = "yanic-temp"
global = "yanic"
#firmware = "yanic-temp"
#model = "yanic-temp"
#autoupdater = "yanic-temp"

# Tagging of the data (optional)
[database.connection.influxdb2.tags] # (6)
# Tags used by Yanic would override the tags from this config
# nodeid, hostname, owner, model, firmware_base, firmware_release,frequency11g and frequency11a are tags which are already used
#tagname1 = "tagvalue 1"
# some useful e.g.:
#system   = "productive"
#site     = "ffhb"
```

1.  Address to connect on InfluxDB2 server.
2.  Token to get acces to InfluxDB2 server.
3.  Set organization using the InfluxDB2 server.

4.  Bucket in which are the data stored.

    Fallback of bucket per measurment, see `[database.connection.influxdb2.buckets]`

5.  Buckets per measurement.

    If not set data `bucket_default` is used.

6.  You could set manuelle tags with inserting into a influxdb.

    Useful if you want to identify the yanic instance when you use multiple own on the same influxdb (e.g. multisites).

    !!! warning
        Tags used by Yanic would override the tags from this config (e.g. `nodeid`, `hostname`, `owner`, `model`, `firmware_base`, `firmware_release`, `frequency11g`, `frequency11a`).

### [[database.connection.graphite]]
Save collected data to a graphite database.

```toml
# Graphite settings
[database.connection.graphite]]
enable   = false
address  = "localhost:2003" # (1)
# Graphite is replacing every "." in the metric name with a slash "/" and uses
# that for the file system hierarchy it generates. it is recommended to at least
# move the metrics out of the root namespace (that would be the empty prefix).
# If you only intend to run one community and only freifunk on your graphite node
# then the prefix can be set to anything (including the empty string) since you
# probably wont care much about "polluting" the namespace.
prefix   = "freifunk" # (2)
```

1.  Address to connect on graphite server.
2.  Graphite is replacing every "." in the metric name with a slash "/" and uses that for the file system hierarchy it generates.
    It is recommended to at least move the metrics out of the root namespace (that would be the empty prefix).
    If you only intend to run one community and only freifunk on your graphite node then the prefix can be set to anything (including the empty string) since you probably wont care much about "polluting" the namespace.

### [[database.connection.respondd]]
Forward collected respondd package to a address
(e.g. to another respondd collector like a central yanic instance or hopglass)

```toml
# respondd (yanic)
# forward collected respondd package to a address
# (e.g. to another respondd collector like a central yanic instance or hopglass)
[[database.connection.respondd]]
enable   = false
# type of network to create a connection
type     = "udp6" # (1)
# destination address to connect/send respondd package
address  = "stats.bremen.freifunk.net:11001" # (2)
```

1.  Type of network to create a connection.

    Known networks are "tcp", "tcp4" (IPv4-only), "tcp6" (IPv6-only), "udp", "udp4" (IPv4-only), "udp6" (IPv6-only), "ip", "ip4" (IPv4-only), "ip6" (IPv6-only), "unix", "unixgram" and "unixpacket".

2.  Destination address to connect/send respondd package.


### [[database.connection.logging]]
This database type is just for, debugging without a real database connection.
A example for other developers for new database types.

```toml
# Logging
[[database.connection.logging]]
enable   = false
path     = "/var/log/yanic.log" # (1)
```

1. Path to file where to store some examples with every line.
