package api

import (
	"encoding/json"
	"fmt"
	"github.com/FreifunkBremen/respond-collector/models"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type ApiAliases struct {
	aliases *models.Aliases
	config  *models.Config
	nodes   *models.Nodes
}

func NewAliases(config *models.Config, router *httprouter.Router, prefix string, nodes *models.Nodes) {
	api := &ApiAliases{
		aliases: models.NewAliases(config),
		nodes:   nodes,
		config:  config,
	}
	router.GET(prefix, api.GetAll)
	router.GET(prefix+"/ansible", api.AnsibleDiff)
	router.GET(prefix+"/cleanup", api.Cleanup)
	router.GET(prefix+"/auth", BasicAuth(api.Cleanup, []byte(config.Webserver.Api.Passphrase)))
	router.GET(prefix+"/alias/:nodeid", api.GetOne)
	router.POST(prefix+"/alias/:nodeid", BasicAuth(api.SaveOne, []byte(config.Webserver.Api.Passphrase)))
}

// clean up the aliases by correct values in nodes
func (api *ApiAliases) cleaner() {
	for key, alias := range api.aliases.List {
		if node := api.nodes.List[key]; node != nil {
			if nodeinfo := node.Nodeinfo; nodeinfo != nil {
				//counter for the diffrent attribute
				count := 3
				if alias.Hostname == nodeinfo.Hostname {
					count -= 1
				}
				if alias.Location.Latitude == nodeinfo.Location.Latitude {
					count -= 1
				}
				if alias.Location.Longtitude == nodeinfo.Location.Longtitude {
					count -= 1
				}
				/*
					if alias.Freq24.TxPower == nodeinfo.Freq24.TxPower {
						count -= 1
					}
					if alias.Freq24.Channel == nodeinfo.Freq24.Channel {
						count -= 1
					}
					if alias.Freq5.TxPower == nodeinfo.Freq5.TxPower {
						count -= 1
					}
					if alias.Freq5.Channel == nodeinfo.Freq5.Channel {
						count -= 1
					}
				*/
				//delete element
				if count <= 0 {
					delete(api.aliases.List, key)
				}
			}
		}
	}
}
func (api *ApiAliases) GetAll(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	jsonOutput(w, r, api.aliases.List)
}

func (api *ApiAliases) GetOne(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if alias := api.aliases.List[ps.ByName("nodeid")]; alias != nil {
		jsonOutput(w, r, alias)
		return
	}
	fmt.Fprint(w, "Not found: ", ps.ByName("nodeid"), "\n")
}

func (api *ApiAliases) SaveOne(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var alias models.Alias

	err := json.NewDecoder(r.Body).Decode(&alias)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Fprint(w, "Decode: ", ps.ByName("nodeid"), "\n")
		return
	}
	api.aliases.Update(ps.ByName("nodeid"), &alias)
	jsonOutput(w, r, alias)
}

func (api *ApiAliases) Cleanup(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	api.cleaner()
	jsonOutput(w, r, api.aliases.List)
}
func (api *ApiAliases) AnsibleDiff(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	api.cleaner()
	jsonOutput(w, r, models.GenerateAnsible(api.nodes, api.aliases.List))
}
