---
layout: default
title: Build and Installation
category_id: 1
permalink: /docs/install.html
---
## go

### Install
```sh
cd /usr/local/
wget https://storage.googleapis.com/golang/go1.9.1.linux-amd64.tar.gz -O go-release-linux-amd64.tar.gz
tar xvf go-release-linux-amd64.tar.gz
rm go-release-linux-amd64.tar.gz
```

### Configure go
Add these lines in your root shell startup file (i.e. `/root/.bashrc`):

```sh
export GOPATH=/opt/go
export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin
```

## Yanic

### Compile
As root:
```sh
go get -v -u github.com/FreifunkBremen/yanic/cmd/...
```

### Install

```bash
cp /opt/go/src/github.com/FreifunkBremen/yanic/contrib/init/linux-systemd/yanic.service /lib/systemd/system/yanic.service
systemctl daemon-reload
```

Before start, you should Configurate the `/etc/yanic.conf`:

```
systemctl start yanic
```

Enable to start on boot:

```
systemctl enable yanic
```
