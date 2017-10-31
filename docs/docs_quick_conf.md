---
layout: default
title: Quick Configuration
category_id: 1
permalink: /docs/quick_configuration.html
---

```sh
cp /opt/go/src/github.com/FreifunkBremen/yanic/config_example.toml /etc/yanic.conf
```

# Quick configuration
For an easy startup you only need to edit the `interfaces` in section
`[respondd]` in file `/etc/yanic.conf`.  

Then create the following folders:
```sh
mkdir -p /var/lib/yanic
mkdir -p /var/www/html/meshviewer/data
```

#### Standalone
If you like to run a standalone meshviewer, just set `enable` in section
`[webserver]` to `true`.

##### Configure the [meshviewer](https://github.com/ffrgb/meshviewer):
set `dataPath` in `config.json` to `/data/` and make the `build` directory
accessible under `/var/www/html/meshviewer`.

#### With webserver (Apache, nginx)
The meshviewer needs the output files like `nodes_path` and `graph_path` inside
the same directory as the `dataPath`. Change the path in the section
`[meshviewer]` accordingly.
