package cmd

import (
	"fmt"
	"os"

	"github.com/FreifunkBremen/yanic/database"
	"github.com/FreifunkBremen/yanic/respond"
	"github.com/FreifunkBremen/yanic/runtime"
	"github.com/FreifunkBremen/yanic/runtime/config"
)

var (
	configPath  string
	collector   *respond.Collector
	connections database.Connection
	nodes       *runtime.Nodes
)

func loadConfig() *config.Config {
	config, err := config.ReadConfigFile(configPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "unable to load config file:", err)
		os.Exit(2)
	}
	return config
}
