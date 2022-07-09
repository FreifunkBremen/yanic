package webserver

import "github.com/FreifunkBremen/yanic/webserver/prometheus"

type Config struct {
	Enable     bool              `toml:"enable"`
	Bind       string            `toml:"bind"`
	Webroot    string            `toml:"webroot"`
	Prometheus prometheus.Config `toml:"prometheus"`
}
