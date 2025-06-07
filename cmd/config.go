package cmd

import (
	"fmt"
	"os"
	"regexp"

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

var envVarRegex = regexp.MustCompile(`\$\{(\w+)\}`)

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

	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	processed := envVarRegex.ReplaceAllStringFunc(string(raw), func(match string) string {
		key := envVarRegex.FindStringSubmatch(match)[1]
		return os.Getenv(key)
	})

	// var config Config
	if _, err := toml.Decode(processed, &config); err != nil {
		return nil, err
	}

	return
}
