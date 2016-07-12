package api

import (
	"encoding/json"
	"fmt"
	"github.com/FreifunkBremen/respond-collector/models"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// 7 nachkommerstellen sollten genug sein (7cm genau)
// http://blog.3960.org/post/7309573249/genauigkeit-bei-geo-koordinaten

const GEOROUND = 0.0000001

func geoEqual(a, b float64) bool {
	if (a-b) < GEOROUND && (b-a) < GEOROUND {
		return true
	}
	return false
}

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
	router.GET(prefix+"/alias/:nodeid", api.GetOne)
	router.POST(prefix+"/alias/:nodeid", BasicAuth(api.SaveOne, []byte(config.Webserver.Api.Passphrase)))
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
	fmt.Print("[api] node updated '", ps.ByName("nodeid"), "'\n")
	jsonOutput(w, r, alias)
}
func (api *ApiAliases) AnsibleDiff(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Print("[api] ansible\n")
	jsonOutput(w, r, models.GenerateAnsible(api.nodes, api.aliases.List))
}
