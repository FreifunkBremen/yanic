# Howto install Yanic

## go

### Install
```sh
cd /usr/local/
wget https://dl.google.com/go/go1.13.1.linux-amd64.tar.gz -O go-release-linux-amd64.tar.gz
tar xvf go-release-linux-amd64.tar.gz
rm go-release-linux-amd64.tar.gz
```

### Configure go
Add these lines in your root shell startup file (e.g. `/root/.bashrc`):
```sh
export GOPATH=/opt/go
export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin
```

## Yanic

### Compile
As root:
```sh
go get -v -u github.com/FreifunkBremen/yanic
```

#### Work with other databases
If you like to use another database solution than influxdb, Pull Requests are
welcome. Just fork this project and create another subpackage within the folder
`database/`. Take this folder as example: `database/logging/`.

### Configure Yanic
```sh
cp /opt/go/src/github.com/FreifunkBremen/yanic/config_example.toml /etc/yanic.conf
```
For an easy startup you only need to edit the `interfaces` in section
`[respondd]` in file `/etc/yanic.conf`.  

Then create the following files and folders:
```sh
mkdir -p /var/lib/yanic
mkdir -p /var/www/html/meshviewer/data
touch /var/log/yanic.log
chown yanic /var/log/yanic.log /var/lib/yanic /var/www/html/meshviewer/data
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
`[[nodes.output.meshviewer]]` accordingly.

### Service
```sh
cp /opt/go/src/github.com/FreifunkBremen/yanic/contrib/init/linux-systemd/yanic.service /lib/systemd/system/yanic.service
systemctl daemon-reload
```

Before start, you should configure yanic by the file `/etc/yanic.conf`:
```sh
systemctl start yanic
```

Enable to start on boot:
```sh
systemctl enable yanic
```

### Update
For an update just stop yanic and then call the same `go` command again (again as root):
```sh
systemctl stop yanic
go get -v -u github.com/FreifunkBremen/yanic
```
Then update the config file, for example look at the diff with the new example:
```sh
diff /opt/go/src/github.com/FreifunkBremen/yanic/config_example.toml /etc/yanic.conf
```
