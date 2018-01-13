package respond

import "github.com/FreifunkBremen/yanic/lib/duration"

type Config struct {
	Enable          bool              `toml:"enable"`
	Synchronize     duration.Duration `toml:"synchronize"`
	Interfaces      []string          `toml:"interfaces"`
	Sites           []string          `toml:"sites"`
	Port            int               `toml:"port"`
	CollectInterval duration.Duration `toml:"collect_interval"`
}
