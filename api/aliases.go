package api

import (
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/FreifunkBremen/respond-collector/models"
)

type ApiAliases struct {
	aliases *models.Aliases
	config		*models.Config
	nodes *models.Nodes
}

func NewAliases (config *models.Config, router *httprouter.Router,prefix string,nodes *models.Nodes) {
	api := &ApiAliases{
		aliases: models.NewAliases(config),
		nodes: nodes,
		config: config,
	}
	router.GET(prefix, api.GetAll)
	router.GET(prefix+"/ansible", api.AnsibleDiff)
	router.GET(prefix+"/alias/:nodeid", api.GetOne)
	router.POST(prefix+"/alias/:nodeid", api.SaveOne)
}

func (api *ApiAliases) GetAll(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	jsonOutput(w,api.aliases.List)
}

func (api *ApiAliases) GetOne(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if alias := api.aliases.List[ps.ByName("nodeid")]; alias !=nil{
		jsonOutput(w,alias)
		return
	}
	fmt.Fprint(w, "Not found: ", ps.ByName("nodeid"),"\n")
}

func (api *ApiAliases) SaveOne(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var alias models.Alias
	err := json.NewDecoder(r.Body).Decode(&alias)

	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Fprint(w, "Decode: ", ps.ByName("nodeid"),"\n")
		return
	}
	api.aliases.Update(ps.ByName("nodeid"),&alias)
	jsonOutput(w,alias)
}

func (api *ApiAliases) AnsibleDiff(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	diff := api.aliases.List
	//TODO diff between List and api.nodes (for run not at all)
	jsonOutput(w,models.GenerateAnsible(api.nodes,diff))
}
