package api

import (
	"encoding/json"
	"fmt"
	"github.com/FreifunkBremen/respond-collector/models"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// GEOROUND : 7 nachkommerstellen sollten genug sein (7cm genau)
// http://blog.3960.org/post/7309573249/genauigkeit-bei-geo-koordinaten
const GEOROUND = 0.0000001

func geoEqual(a, b float64) bool {
	if (a-b) < GEOROUND && (b-a) < GEOROUND {
		return true
	}
	return false
}

// AliasesAPI struct for API
type AliasesAPI struct {
	aliases *models.Aliases
	config  *models.Config
	nodes   *models.Nodes
}

// NewAliases Bind to API
func NewAliases(config *models.Config, router *httprouter.Router, prefix string, nodes *models.Nodes) {
	api := &AliasesAPI{
		aliases: models.NewAliases(config),
		nodes:   nodes,
		config:  config,
	}
	router.GET(prefix, api.GetAll)
	router.GET(prefix+"/ansible", api.Ansible)
	router.GET(prefix+"/alias/:nodeid", api.GetOne)
	router.POST(prefix+"/alias/:nodeid", BasicAuth(api.SaveOne, []byte(config.Webserver.API.Passphrase)))
}

// GetAll request for get all aliases
func (api *AliasesAPI) GetAll(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	jsonOutput(w, r, api.aliases.List)
}

// GetOne request for get one alias
func (api *AliasesAPI) GetOne(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if alias := api.aliases.List[ps.ByName("nodeid")]; alias != nil {
		jsonOutput(w, r, alias)
		return
	}
	fmt.Fprint(w, "Not found: ", ps.ByName("nodeid"), "\n")
}

// SaveOne request for save a alias
func (api *AliasesAPI) SaveOne(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

// Ansible json output
func (api *AliasesAPI) Ansible(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Print("[api] ansible\n")
	jsonOutput(w, r, models.GenerateAnsible(api.nodes, api.aliases.List))
}
