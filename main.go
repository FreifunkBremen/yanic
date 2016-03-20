package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/FreifunkBremen/respond-collector/data"
	"github.com/FreifunkBremen/respond-collector/models"
	"github.com/FreifunkBremen/respond-collector/respond"
	"github.com/FreifunkBremen/respond-collector/websocketserver"
)

var (
	configFile       string
	config           *models.Config
	wsserverForNodes *websocketserver.Server
	collector        *respond.Collector
	statsDb          *StatsDb
	nodes            = models.NewNodes()
	//aliases          = models.NewNodes()
)

func main() {
	flag.StringVar(&configFile, "config", "config.yml", "path of configuration file (default:config.yaml)")
	flag.Parse()
	config = models.ConfigReadFile(configFile)

	if config.Nodes.Enable {
		go nodes.Saver(config)
	}

	if config.Webserver.Enable {
		if config.Webserver.WebsocketNode {
			wsserverForNodes = websocketserver.NewServer("/nodes")
			go wsserverForNodes.Listen()
		}
		http.Handle("/", http.FileServer(http.Dir(config.Webserver.Webroot)))
	}

	if config.Influxdb.Enable {
		statsDb = NewStatsDb()
	}

	if config.Respondd.Enable {
		collectInterval := time.Second * time.Duration(config.Respondd.CollectInterval)
		collector = respond.NewCollector("nodeinfo statistics neighbours", collectInterval, onReceive)
	}

	// TODO bad
	if config.Webserver.Enable {
		log.Fatal(http.ListenAndServe(net.JoinHostPort(config.Webserver.Address, config.Webserver.Port), nil))
	}

	// Wait for INT/TERM
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigs
	log.Println("received", sig)

	// Close everything at the end
	if wsserverForNodes != nil {
		wsserverForNodes.Close()
	}
	if collector != nil {
		collector.Close()
	}
	if statsDb != nil {
		statsDb.Close()
	}
}

// called for every parsed announced-message
func onReceive(addr net.UDPAddr, res *data.ResponseData) {

	// Search for NodeID
	var nodeId string
	if val := res.NodeInfo; val != nil {
		nodeId = val.NodeId
	} else if val := res.Neighbours; val != nil {
		nodeId = val.NodeId
	} else if val := res.Statistics; val != nil {
		nodeId = val.NodeId
	}

	// Updates nodes if NodeID found
	if len(nodeId) != 12 {
		log.Printf("invalid NodeID '%s' from %s", nodeId, addr.String())
		return
	}

	nodes.Update(nodeId, res)

	if val := res.Statistics; val != nil && statsDb != nil {
		statsDb.Add(val)
	}
}
