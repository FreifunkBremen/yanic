package models

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sync"
	"time"

	"github.com/FreifunkBremen/respond-collector/data"
)

// Alias a change request for other nodes
type Alias struct {
	Hostname string         `json:"hostname,omitempty"`
	Location *data.Location `json:"location,omitempty"`
	Wireless *data.Wireless `json:"wireless,omitempty"`
	Owner    string         `json:"owner,omitempty"`
}

// Aliases struct: cache DB of Node's structs
type Aliases struct {
	List   map[string]*Alias `json:"nodes"` // the current nodemap, indexed by node ID
	config *Config
	sync.Mutex
}

// NewAliases create Nodes structs
func NewAliases(config *Config) *Aliases {
	aliases := &Aliases{
		List:   make(map[string]*Alias),
		config: config,
	}

	if config.Nodes.AliasesPath != "" {
		aliases.load()
	}
	go aliases.worker()

	return aliases
}

// Update a alias in aliases cache
func (e *Aliases) Update(nodeID string, newalias *Alias) {
	e.Lock()
	e.List[nodeID] = newalias
	e.Unlock()

}

func (e *Aliases) load() {
	path := e.config.Nodes.AliasesPath
	log.Println("loading", path)

	if data, err := ioutil.ReadFile(path); err == nil {
		if err = json.Unmarshal(data, e); err == nil {
			log.Println("loaded", len(e.List), "aliases")
		} else {
			log.Println("failed to unmarshal nodes:", err)
		}

	} else {
		log.Println("failed loading cached nodes:", err)
	}
}

// Periodically saves the cached DB to json file
func (e *Aliases) worker() {
	c := time.Tick(time.Second * 5)

	for range c {
		log.Println("saving", len(e.List), "aliases")
		e.Lock()
		save(e, e.config.Nodes.AliasesPath)
		e.Unlock()
	}
}
