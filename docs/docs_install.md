---
layout: default
title: Build and Installation
category_id: 1
permalink: /docs/install.html
---

## Install go

```sh
cd /usr/local/
wget https://storage.googleapis.com/golang/go1.9.1.linux-amd64.tar.gz -O go-release-linux-amd64.tar.gz
tar xvf go-release-linux-amd64.tar.gz
rm go-release-linux-amd64.tar.gz
```

### Configurate

put this lines into a shell place at root for easy yanic install:

```sh
export GOPATH=/opt/go
export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin
```

or put this lines also into a shell place to use go by normal user:

```sh
export GOPATH=~/go
export PATH=$PATH:$GOPATH/bin
```

## Yanic

### Build

```sh
go get -v -u github.com/FreifunkBremen/yanic/cmd/...
```

### Install

```bash
cp /opt/go/src/github.com/FreifunkBremen/yanic/contrib/init/linux-systemd/yanic.service /lib/systemd/system/
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
