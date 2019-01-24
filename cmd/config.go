package cmd

import (
	"io/ioutil"

	"github.com/naoina/toml"

	"github.com/FreifunkBremen/yanic/respond"
	"github.com/FreifunkBremen/yanic/runtime"
)

var (
	configPath string
	collector  *respond.Collector
	nodes      *runtime.Nodes
)

// ReadConfigFile reads a config model from path of a yml file
func ReadConfigFile(path string, config interface{}) error {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = toml.Unmarshal(file, config)
	if err != nil {
		return err
	}

	return nil
}
