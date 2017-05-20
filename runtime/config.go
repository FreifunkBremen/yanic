package runtime

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

//Config the config File of this daemon
type Config struct {
	Respondd struct {
		Enable          bool     `toml:"enable"`
		Synchronize     Duration `toml:"synchronize"`
		Interface       string   `toml:"interface"`
		Port            int      `toml:"port"`
		CollectInterval Duration `toml:"collect_interval"`
	}
	Webserver struct {
		Enable  bool   `toml:"enable"`
		Bind    string `toml:"bind"`
		Webroot string `toml:"webroot"`
	}
	Nodes struct {
		Enable       bool     `toml:"enable"`
		StatePath    string   `toml:"state_path"`
		SaveInterval Duration `toml:"save_interval"` // Save nodes periodically
		OfflineAfter Duration `toml:"offline_after"` // Set node to offline if not seen within this period
		PruneAfter   Duration `toml:"prune_after"`   // Remove nodes after n days of inactivity
		Output       map[string][]interface{}
	}
	Meshviewer struct {
		Version   int    `toml:"version"`
		NodesPath string `toml:"nodes_path"`
		GraphPath string `toml:"graph_path"`
	}
	Database struct {
		DeleteInterval Duration `toml:"delete_interval"` // Delete stats of nodes every n minutes
		DeleteAfter    Duration `toml:"delete_after"`    // Delete stats of nodes till now-deletetill n minutes
		Connection     map[string][]interface{}
	}
}

// ReadConfigFile reads a config model from path of a yml file
func ReadConfigFile(path string) (config *Config, err error) {
	config = &Config{}

	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	err = toml.Unmarshal(file, config)
	if err != nil {
		panic(err)
	}

	return
}
