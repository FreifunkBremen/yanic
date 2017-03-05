# Howto install Yanic

## go
### Install
```sh
cd /usr/local/
wget https://storage.googleapis.com/golang/go1.8.linux-amd64.tar.gz
tar xvf go1.8.linux-amd64.tar.gz
rm go1.8.linux-amd64.tar.gz
```

### Configurate
put this lines into a shell place at root:
```sh
export GOPATH=/opt/go
export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin
```

put this lines also into a shell place to use go by normal user:
```sh
export GOPATH=~/go
export PATH=$PATH:$GOPATH/bin
```

## Yanic

### Compile
```sh
go get -v -u github.com/FreifunkBremen/yanic/cmd/...
```

### Configurate
```sh
cp /opt/go/src/github.com/FreifunkBremen/yanic/config_example.toml /etc/yanic.conf
```
You only need to edit `/etc/yanic.conf` under section `[respondd]` the `interface` for a easy startup.
And create the following folders:
```sh
mkdir -p /var/lib/collector
mkdir -p /var/www/html/meshviewer/data
```

#### Standalone
If you like to run a meshviewer standalone, just set `enable` under section `[webserver]` to `true`.
Configurate the [meshviewer](https://github.com/ffrgb/meshviewer) set `dataPath` in `config.json` to `/data/` and put the `build` directory under `/var/www/html/meshviewer`.

#### With webserver (Apache, nginx)
Change following path under section `[nodes]` to what you need.
For `nodes_path` and `graph_path` should be under the same folder for a meshviewer.

### Service
```bash
cp /opt/go/src/github.com/FreifunkBremen/yanic/init/linux-systemd/yanic.service /lib/systemd/systemd
systemctl daemon-reload
systemctl start yanic
systemctl enable yanic
```
