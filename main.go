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
	"github.com/FreifunkBremen/respond-collector/database"
	"github.com/FreifunkBremen/respond-collector/models"
	"github.com/FreifunkBremen/respond-collector/respond"
	"github.com/FreifunkBremen/respond-collector/rrd"
)

var (
	configFile string
	config     *models.Config
	collector  *respond.Collector
	db         *database.DB
	nodes      *models.Nodes
)

func main() {
	var importPath string
	flag.StringVar(&importPath, "import", "", "import global statistics from the given RRD file, requires influxdb")
	flag.StringVar(&configFile, "config", "config.yml", "path of configuration file (default:config.yaml)")
	flag.Parse()
	config = models.ReadConfigFile(configFile)

	if config.Influxdb.Enable {
		db = database.New(config)
		defer db.Close()

		if importPath != "" {
			importRRD(importPath)
			return
		}
	}

	nodes = models.NewNodes(config)
	if config.Respondd.Enable {
		collectInterval := time.Second * time.Duration(config.Respondd.CollectInterval)
		collector = respond.NewCollector(db, nodes, collectInterval, config.Respondd.Interface)
		defer collector.Close()
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
}

func importRRD(path string) {
	log.Println("importing RRD from", path)
	for ds := range rrd.Read(path) {
		db.AddPoint(
			database.MeasurementGlobal,
			nil,
			map[string]interface{}{
				"nodes":         ds.Nodes,
				"clients.total": ds.Clients,
			},
			ds.Time,
		)
	}
}
