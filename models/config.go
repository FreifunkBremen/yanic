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
		Api     struct {
			Passphrase string `yaml:"passphrase"`
			NewNodes   bool   `yaml:"newnodes"`
			Aliases    bool   `yaml:"aliases"`
		} `yaml:"api"`
	} `yaml:"webserver"`
	Nodes struct {
		Enable        bool   `yaml:"enable"`
		NodesPath     string `yaml:"nodes_path"`
		NodesMiniPath string `yaml:"nodesmini_path"`
		GraphsPath    string `yaml:"graphs_path"`
		AliasesPath   string `yaml:"aliases_path"`
		SaveInterval  int    `yaml:"saveinterval"`
	} `yaml:"nodes"`
	Influxdb struct {
		Enable   bool   `yaml:"enable"`
		Addr     string `yaml:"host"`
		Database string `yaml:"database"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	}
}

// reads a config models by path to a yml file
func ReadConfigFile(path string) *Config {
	config := &Config{}
	file, _ := ioutil.ReadFile(path)
	err := yaml.Unmarshal(file, &config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}
