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

	"github.com/ffdo/node-informant/gluon-collector/data"
	"github.com/monitormap/micro-daemon/models"
	"github.com/monitormap/micro-daemon/respond"
	"github.com/monitormap/micro-daemon/websocketserver"
)

var (
	configFile       string
	config           *models.Config
	wsserverForNodes *websocketserver.Server
	multiCollector   *respond.MultiCollector
	statsDb          *StatsDb
	nodes            = models.NewNodes()
	//aliases          = models.NewNodes()
)

func main() {
	flag.StringVar(&configFile, "c", "config.yml", "path of configuration file (default:config.yaml)")
	flag.Parse()
	config = models.ConfigReadFile(configFile)

	collectInterval := time.Second * time.Duration(config.Respondd.CollectInterval)

	if config.Nodes.Enable {
		go nodes.Saver(config)
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

	if config.Influxdb.Enable {
		statsDb = NewStatsDb()
	}

	if config.Respondd.Enable {
		multiCollector = respond.NewMultiCollector(collectInterval, func(addr net.UDPAddr, msg interface{}) {
			switch msg := msg.(type) {
			case *data.NodeInfo:
				nodes.Get(msg.NodeId).Nodeinfo = msg
			case *data.NeighbourStruct:
				nodes.Get(msg.NodeId).Neighbours = msg
			case *data.StatisticsStruct:
				nodes.Get(msg.NodeId).Statistics = msg
				if statsDb != nil {
					statsDb.Add(msg)
				}
			default:
				log.Println("unknown message:", msg)
			}
		})
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
	if wsserverForNodes != nil {
		wsserverForNodes.Close()
	}
	if multiCollector != nil {
		multiCollector.Close()
	}
	if statsDb != nil {
		statsDb.Close()
	}
}
