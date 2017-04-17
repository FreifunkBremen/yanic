package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/FreifunkBremen/yanic/database"
	"github.com/FreifunkBremen/yanic/database/all"
	"github.com/FreifunkBremen/yanic/meshviewer"
	"github.com/FreifunkBremen/yanic/respond"
	"github.com/FreifunkBremen/yanic/rrd"
	"github.com/FreifunkBremen/yanic/runtime"
	"github.com/FreifunkBremen/yanic/webserver"
)

var (
	configFile  string
	config      *runtime.Config
	collector   *respond.Collector
	connections database.Connection
	nodes       *runtime.Nodes
)

func main() {
	var importPath string
	var timestamps bool
	flag.StringVar(&importPath, "import", "", "import global statistics from the given RRD file, requires influxdb")
	flag.StringVar(&configFile, "config", "config.toml", "path of configuration file (default:config.yaml)")
	flag.BoolVar(&timestamps, "timestamps", true, "print timestamps in output")
	flag.Parse()

	if !timestamps {
		log.SetFlags(0)
	}
	log.Println("Yanic say hello")

	config, err := runtime.ReadConfigFile(configFile)
	if err != nil {
		panic(err)
	}

	connections, err = all.Connect(config.Database.Connection)
	if err != nil {
		panic(err)
	}
	database.Start(connections, config)
	defer database.Close(connections)

	if connections != nil && importPath != "" {
		importRRD(importPath)
		return
	}

	nodes = runtime.NewNodes(config)
	nodes.Start()
	meshviewer.Start(config, nodes)

	if config.Respondd.Enable {
		collector = respond.NewCollector(connections, nodes, config.Respondd.Interface, config.Respondd.Port)
		collector.Start(config.Respondd.CollectInterval.Duration)
		defer collector.Close()
	}

	if config.Webserver.Enable {
		log.Println("starting webserver on", config.Webserver.Bind)
		srv := webserver.New(config.Webserver.Bind, config.Webserver.Webroot)
		go srv.Close()
	}

	// Wait for INT/TERM
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigs
	log.Println("received", sig)
}

func importRRD(path string) {
	log.Println("importing RRD from", path)
	for ds := range rrd.Read(path) {
		connections.InsertGlobals(
			&runtime.GlobalStats{
				Nodes:   uint32(ds.Nodes),
				Clients: uint32(ds.Clients),
			},
			ds.Time,
		)
	}
}
