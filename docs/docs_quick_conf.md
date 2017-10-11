---
layout: default
title: Quick Configuration
category_id: 1
permalink: /docs/quick_configuration.html
---

```bash
cp /opt/go/src/github.com/FreifunkBremen/yanic/config_example.toml /etc/yanic.conf
```



# Quick configuration



You only need to edit `/etc/yanic.conf` under section `[respondd]` the `interface` for a easy startup.

And create the following folders:



```bash
mkdir -p /var/lib/yanic
mkdir -p /var/www/html/meshviewer/data
```

## Standalone

If you like to run a meshviewer standalone, just set `enable` under section `[webserver]` to `true`.

Configurate the [meshviewer](https://github.com/ffrgb/meshviewer) set `dataPath` in `config.json` to `/data/` and put the `build` directory under `/var/www/html/meshviewer`.



## With webserver \(Apache, nginx\)

Change following path under section `[nodes]` to what you need.

For `nodes_path` and `graph_path` should be under the same folder for a meshviewer.
