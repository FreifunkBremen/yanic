# Build and Installation

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

### Install

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
