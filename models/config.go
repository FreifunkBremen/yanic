package models

import (
	"io/ioutil"

	"github.com/influxdata/toml"
)

//Config the config File of this daemon
type Config struct {
	Respondd struct {
		Enable          bool
		Interface       string
		CollectInterval Duration
	}
	Webserver struct {
		Enable  bool
		Bind    string
		Webroot string
	}
	Nodes struct {
		Enable       bool
		NodesVersion int
		NodesPath    string
		GraphPath    string
		StatePath    string
		SaveInterval Duration // Save nodes periodically
		PruneAfter   Duration // Remove nodes after n days of inactivity
	}
	Influxdb struct {
		Enable         bool
		Address        string
		Database       string
		Username       string
		Password       string
		SaveInterval   Duration // Save nodes every n seconds
		DeleteInterval Duration // Delete stats of nodes every n minutes
		DeleteAfter    Duration // Delete stats of nodes till now-deletetill n minutes
	}
}

// ReadConfigFile reads a config model from path of a yml file
func ReadConfigFile(path string) *Config {
	config := &Config{}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	if err := toml.Unmarshal(file, config); err != nil {
		panic(err)
	}

	return config
}
