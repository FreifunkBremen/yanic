package config

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"github.com/FreifunkBremen/yanic/database"
	"github.com/FreifunkBremen/yanic/lib/duration"
	"github.com/FreifunkBremen/yanic/runtime"
)

//Config the config File of this daemon
type Config struct {
	Respondd struct {
		Enable          bool              `toml:"enable"`
		Synchronize     duration.Duration `toml:"synchronize"`
		Interfaces      []string          `toml:"interfaces"`
		Sites           []string          `toml:"sites"`
		Port            int               `toml:"port"`
		CollectInterval duration.Duration `toml:"collect_interval"`
	}
	Webserver struct {
		Enable  bool   `toml:"enable"`
		Bind    string `toml:"bind"`
		Webroot string `toml:"webroot"`
	}
	Nodes    runtime.NodesConfig
	Database database.Config
}

// ReadConfigFile reads a config model from path of a yml file
func ReadConfigFile(path string) (config *Config, err error) {
	config = &Config{}

	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = toml.Unmarshal(file, config)
	if err != nil {
		return nil, err
	}

	return
}
