package models

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

//Config the config File of this daemon
type Config struct {
	Respondd struct {
		Enable          bool   `yaml:"enable"`
		Interface       string `yaml:"interface"`
		CollectInterval int    `yaml:"collectinterval"`
	} `yaml:"respondd"`
	Webserver struct {
		Enable  bool   `yaml:"enable"`
		Port    string `yaml:"port"`
		Address string `yaml:"address"`
		Webroot string `yaml:"webroot"`
		API     struct {
			Passphrase string `yaml:"passphrase"`
			NewNodes   bool   `yaml:"newnodes"`
			Aliases    bool   `yaml:"aliases"`
		} `yaml:"api"`
	} `yaml:"webserver"`
	Nodes struct {
		Enable           bool   `yaml:"enable"`
		NodesDynamicPath string `yaml:"nodes_path"`
		NodesV1Path      string `yaml:"nodesv1_path"`
		NodesV2Path      string `yaml:"nodesv2_path"`
		GraphsPath       string `yaml:"graphs_path"`
		AliasesPath      string `yaml:"aliases_path"`
		SaveInterval     int    `yaml:"saveinterval"` // Save nodes every n seconds
		MaxAge           int    `yaml:"max_age"`      // Remove nodes after n days of inactivity
	} `yaml:"nodes"`
	Influxdb struct {
		Enable         bool   `yaml:"enable"`
		Addr           string `yaml:"host"`
		Database       string `yaml:"database"`
		Username       string `yaml:"username"`
		Password       string `yaml:"password"`
		SaveInterval   int    `yaml:"saveinterval"`   // Save nodes every n seconds
		DeleteInterval int    `yaml:"deleteinterval"` // Delete stats of nodes every n minutes
		DeleteTill     int    `yaml:"deletetill"`     // Delete stats of nodes till now-deletetill n minutes
	}
}

// ReadConfigFile reads a config model from path of a yml file
func ReadConfigFile(path string) *Config {
	config := &Config{}
	file, _ := ioutil.ReadFile(path)
	err := yaml.Unmarshal(file, &config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}
