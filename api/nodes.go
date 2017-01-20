package api

import (
	"github.com/FreifunkBremen/respond-collector/models"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// NodesAPI struct for API
type NodesAPI struct {
	config *models.Config
	nodes  *models.Nodes
}

// NewNodes Bind to API
func NewNodes(config *models.Config, router *httprouter.Router, prefix string, nodes *models.Nodes) {
	api := &NodesAPI{
		nodes:  nodes,
		config: config,
	}
	router.GET(prefix, api.GetAll)
}

// GetAll request for get all nodes
func (api *NodesAPI) GetAll(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	jsonOutput(w, r, api.nodes.List)
}
