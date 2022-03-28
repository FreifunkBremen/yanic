package cmd

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"

	"github.com/FreifunkBremen/yanic/database"
	"github.com/FreifunkBremen/yanic/respond"
	"github.com/FreifunkBremen/yanic/runtime"
	"github.com/FreifunkBremen/yanic/webserver"
)

// Config represents the whole configuration
type Config struct {
	Respondd  respond.Config
	Webserver webserver.Config
	Nodes     runtime.NodesConfig
	Database  database.Config
}

var (
	configPath string
	collector  *respond.Collector
	nodes      *runtime.Nodes
)

func loadConfig() *Config {
	config, err := ReadConfigFile(configPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "unable to load config file:", err)
		os.Exit(2)
	}
	return config
}

// ReadConfigFile reads a config model from path of a yml file
func ReadConfigFile(path string) (config *Config, err error) {
	config = &Config{}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	_, err = toml.NewDecoder(file).Decode(config)
	if err != nil {
		return nil, err
	}

	return
}
