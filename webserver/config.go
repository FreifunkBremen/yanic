package webserver

type Config struct {
	Enable     bool        `toml:"enable"`
	Bind       string      `toml:"bind"`
	Webroot    string      `toml:"webroot"`
	Prometheus *Prometheus `toml:"prometheus"`
}
