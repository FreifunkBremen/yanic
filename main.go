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

	"github.com/NYTimes/gziphandler"
	"github.com/julienschmidt/httprouter"

	"github.com/FreifunkBremen/respond-collector/api"
	"github.com/FreifunkBremen/respond-collector/data"
	"github.com/FreifunkBremen/respond-collector/models"
	"github.com/FreifunkBremen/respond-collector/respond"
)

var (
	configFile string
	config     *models.Config
	collector  *respond.Collector
	statsDb    *StatsDb
	nodes      *models.Nodes
)

func main() {
	flag.StringVar(&configFile, "config", "config.yml", "path of configuration file (default:config.yaml)")
	flag.Parse()
	config = models.ReadConfigFile(configFile)
	nodes = models.NewNodes(config)

	if config.Influxdb.Enable {
		statsDb = NewStatsDb()
	}

	if config.Respondd.Enable {
		collectInterval := time.Second * time.Duration(config.Respondd.CollectInterval)
		collector = respond.NewCollector("nodeinfo statistics neighbours", collectInterval, onReceive, config.Respondd.Interface)
	}

	if config.Webserver.Enable {
		router := httprouter.New()
		if config.Webserver.Api.NewNodes {
			api.NewNodes(config, router, "/api/nodes", nodes)
			log.Println("api nodes started")
		}
		if config.Webserver.Api.Aliases {
			api.NewAliases(config, router, "/api/aliases", nodes)
			log.Println("api aliases started")
		}
		router.NotFound = gziphandler.GzipHandler(http.FileServer(http.Dir(config.Webserver.Webroot)))

		address := net.JoinHostPort(config.Webserver.Address, config.Webserver.Port)
		log.Println("starting webserver on", address)
		// TODO bad
		log.Fatal(http.ListenAndServe(address, router))
	}

	// Wait for INT/TERM
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigs
	log.Println("received", sig)

	// Close everything at the end
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
		if node := res.NodeInfo; val != nil && statsDb != nil {
			statsDb.Add(val, node)
		} else {
			statsDb.Add(val, nil)
		}
	}
}
