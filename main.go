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
	wsserverForNodes  = websocketserver.NewServer("/nodes")
	respondDaemon     *respond.Daemon
	nodes             = models.NewNodes()
	aliases           = models.NewNodes()
	listenAddr        string
	listenPort        string
	collectInterval   time.Duration
	httpDir           string
	outputNodesFile   string
	outputAliasesFile string
	saveInterval      time.Duration
)

func main() {
	var collectSeconds, saveSeconds int

	flag.StringVar(&listenAddr, "host", "", "path aliases.json file")
	flag.StringVar(&listenPort, "port", "8080", "path aliases.json file")
	flag.IntVar(&collectSeconds, "collectInterval", 15, "interval for data collections")
	flag.StringVar(&httpDir, "httpdir", "webroot", "a implemented static file webserver")
	flag.StringVar(&outputNodesFile, "path-nodes", "webroot/nodes.json", "path nodes.json file")
	flag.StringVar(&outputAliasesFile, "path-aliases", "webroot/aliases.json", "path aliases.json file")
	flag.IntVar(&saveSeconds, "saveInterval", 60, "interval for data saving")
	flag.Parse()

	collectInterval = time.Second * time.Duration(collectSeconds)
	saveInterval = time.Second * time.Duration(saveSeconds)

	go wsserverForNodes.Listen()
	go nodes.Saver(outputNodesFile, saveInterval)
	go aliases.Saver(outputAliasesFile, saveInterval)
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

	http.Handle("/", http.FileServer(http.Dir(httpDir)))
	//TODO bad
	log.Fatal(http.ListenAndServe(net.JoinHostPort(listenAddr, listenPort), nil))

	// Wait for End
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigs
	log.Println("received", sig)

	// Close everything at the end
	wsserverForNodes.Close()
	respondDaemon.Close()
}
