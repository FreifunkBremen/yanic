package main

import (
	"encoding/json"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ffdo/node-informant/gluon-collector/data"
	"github.com/monitormap/micro-daemon/models"
	"github.com/monitormap/micro-daemon/respond"
	"github.com/monitormap/micro-daemon/websocketserver"
)

var (
	configFile       string
	config           *models.Config
	wsserverForNodes *websocketserver.Server
	respondDaemon    *respond.Daemon
	nodes            = models.NewNodes()
	//aliases          = models.NewNodes()
)

func main() {
	flag.StringVar(&configFile, "c", "config.yml", "path of configuration file (default:config.yaml)")
	flag.Parse()
	config = models.ConfigReadFile(configFile)

	collectInterval := time.Second * time.Duration(config.Responedd.CollectInterval)
	saveInterval := time.Second * time.Duration(config.Nodes.SaveInterval)

	if config.Nodes.Enable {
		go nodes.Saver(config.Nodes.NodesPath, config.Nodes.GraphsPath, saveInterval)
	}
	if config.Nodes.AliasesEnable {
		// FIXME what does this do?
		//go aliases.Saver(config.Nodes.AliasesPath, saveInterval)
	}

	if config.Webserver.Enable {
		if config.Webserver.WebsocketNode {
			wsserverForNodes = websocketserver.NewServer("/nodes")
			go wsserverForNodes.Listen()
		}
		http.Handle("/", http.FileServer(http.Dir(config.Webserver.Webroot)))
	}

	if config.Responedd.Enable {
		respondDaemon = respond.NewDaemon(func(coll *respond.Collector, res *respond.Response) {

			switch coll.CollectType {
			case "neighbours":
				result := &data.NeighbourStruct{}
				if json.Unmarshal(res.Raw, result) == nil {
					node := nodes.Get(result.NodeId)
					node.Neighbours = result
				}
			case "nodeinfo":
				result := &data.NodeInfo{}
				if json.Unmarshal(res.Raw, result) == nil {
					node := nodes.Get(result.NodeId)
					node.Nodeinfo = result
				}
			case "statistics":
				result := &data.StatisticsStruct{}
				if json.Unmarshal(res.Raw, result) == nil {
					node := nodes.Get(result.NodeId)
					node.Statistics = result
				}
			default:
				log.Println("unknown CollectType:", coll.CollectType)
			}
		})
		go respondDaemon.ListenAndSend(collectInterval)
	}

	//TODO bad
	if config.Webserver.Enable {
		log.Fatal(http.ListenAndServe(net.JoinHostPort(config.Webserver.Address, config.Webserver.Port), nil))
	}
	// Wait for End
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigs
	log.Println("received", sig)

	// Close everything at the end
	wsserverForNodes.Close()
	respondDaemon.Close()
}
