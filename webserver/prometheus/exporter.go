package prometheus

import (
	"io"
	"net"
	"net/http"
	"time"

	"github.com/bdlm/log"

	"github.com/FreifunkBremen/yanic/respond"
	"github.com/FreifunkBremen/yanic/runtime"
)

type Exporter struct {
	config Config
	srv *http.Server
	coll *respond.Collector
	nodes *runtime.Nodes
}


func (ex *Exporter) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var ip net.IP
	nodeID := ""

	queryValues := req.URL.Query()

	if nodeIDs := queryValues["node_id"]; len(nodeIDs) > 0 {
		nodeID = nodeIDs[0]
		node, ok := ex.nodes.List[nodeID]
		if !ok || node.Address == nil {
			http.Error(res, "not able to get node by cached nodeid", http.StatusNotFound)
			return
		}
		ip = node.Address.IP
		if ex.writeNode(res, node) {
			log.WithFields(map[string]interface{}{
				"ip": ip,
				"node_id": nodeID,
			}).Debug("take node from cache")
			return
		}
	} else if ipstr := queryValues["ip"]; len(ipstr) > 0 {
		ip = net.ParseIP(ipstr[0])
		if ip == nil {
			http.Error(res, "not able to parse ip address", http.StatusBadRequest)
			return
		}
		node_select := ex.nodes.Select(func(n *runtime.Node) bool {
			n_addr := n.Address
			nodeID = n.Nodeinfo.NodeID
			return n_addr != nil && ip.Equal(n_addr.IP)
		})
		getOne := false
		if len(node_select) == 1 {
			if ex.writeNode(res, node_select[0]) {
				getOne = true
			}
		} else if len(node_select) > 1 {
			log.Error("strange count of nodes")
		}

		if getOne {
			log.WithFields(map[string]interface{}{
				"ip": ip,
				"node_id": nodeID,
			}).Debug("take node from cache")
			return
		}
	} else {
		http.Error(res, "please request with ?ip= or ?node_id=", http.StatusNotFound)
		return
	}

	// send request
	ex.coll.SendPacket(ip)

	// wait
	log.WithFields(map[string]interface{}{
		"ip": ip,
		"node_id": nodeID,
	}).Debug("waited for")
	time.Sleep(ex.config.Wait.Duration)

	// result
	node, ok := ex.nodes.List[nodeID]
	if !ok {
		http.Error(res, "not able to fetch this node", http.StatusGatewayTimeout)
		return
	}
	if ex.writeNode(res, node) {
		return
	}
	http.Error(res, "not able to fetch new values from this node", http.StatusGatewayTimeout)
}

func (ex *Exporter) writeNode(res http.ResponseWriter, node *runtime.Node) bool {
	if !time.Now().Before(node.Lastseen.GetTime().Add(ex.config.Outdated.Duration)) {
		return false
	}
	metrics := MetricsFromNode(ex.nodes, node)
	for _, m := range metrics {
		str, err := m.String()
		if err == nil {
			io.WriteString(res, str + "\n")
		} else {
			logger := log.WithField("database", "prometheus")
			if nodeinfo := node.Nodeinfo; nodeinfo != nil {
				logger = logger.WithField("node_id", nodeinfo.NodeID)
			}
			logger.Warnf("not able to get metrics from node: %s", err)
			http.Error(res, "not able to generate metric from node", http.StatusInternalServerError)
		}
	}

	return true
}
