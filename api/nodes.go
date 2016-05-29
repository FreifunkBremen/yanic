package api

import (
	"github.com/FreifunkBremen/respond-collector/models"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type ApiNodes struct {
	config *models.Config
	nodes  *models.Nodes
}

func NewNodes(config *models.Config, router *httprouter.Router, prefix string, nodes *models.Nodes) {
	api := &ApiNodes{
		nodes:  nodes,
		config: config,
	}
	router.GET(prefix, api.GetAll)
}

func (api *ApiNodes) GetAll(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	jsonOutput(w, r, api.nodes.List)
}
