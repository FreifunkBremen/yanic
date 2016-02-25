package main

import (
	"encoding/json"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"syscall"
	"time"

	"github.com/monitormap/micro-daemon/models"
	"github.com/monitormap/micro-daemon/responed"
	"github.com/monitormap/micro-daemon/websocketserver"
)

var (
	wsserverForNodes  = websocketserver.NewServer("/nodes")
	responedDaemon    *responed.Daemon
	nodes             = models.NewNodes()
	aliases           = models.NewNodes()
	outputNodesFile   string
	outputAliasesFile string
	collectInterval   time.Duration
	saveInterval      time.Duration
	listenPort        string
	listenAddr        string
)

func main() {
	var collectSeconds, saveSeconds int

	flag.StringVar(&outputNodesFile, "output", "webroot/nodes.json", "path nodes.json file")
	flag.StringVar(&outputAliasesFile, "aliases", "webroot/aliases.json", "path aliases.json file")
	flag.StringVar(&listenPort, "p", "8080", "path aliases.json file")
	flag.StringVar(&listenAddr, "h", "", "path aliases.json file")
	flag.IntVar(&saveSeconds, "saveInterval", 60, "interval for data saving")
	flag.IntVar(&collectSeconds, "collectInterval", 15, "interval for data collections")
	flag.Parse()

	collectInterval = time.Second * time.Duration(collectSeconds)
	saveInterval = time.Second * time.Duration(saveSeconds)

	go wsserverForNodes.Listen()
	go nodes.Saver(outputNodesFile, saveInterval)
	go aliases.Saver(outputAliasesFile, saveInterval)
	responedDaemon = responed.NewDaemon(func(coll *responed.Collector, res *responed.Response) {
		var result map[string]interface{}
		json.Unmarshal(res.Raw, &result)

		nodeID, _ := result["node_id"].(string)

		if nodeID == "" {
			log.Println("unable to parse node_id")
			return
		}

		node := nodes.Get(nodeID)

		// Set result
		elem := reflect.ValueOf(node).Elem()
		field := elem.FieldByName(strings.Title(coll.CollectType))

		if !reflect.DeepEqual(field, result) {
			wsserverForNodes.SendAll(node)
		}

		field.Set(reflect.ValueOf(result))
	})
	go responedDaemon.ListenAndSend(collectInterval)

	http.Handle("/", http.FileServer(http.Dir("webroot")))
	//TODO bad
	log.Fatal(http.ListenAndServe(net.JoinHostPort(listenAddr, listenPort), nil))

	// Wait for End
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigs
	log.Println("received", sig)

	// Close everything at the end
	wsserverForNodes.Close()
	responedDaemon.Close()
}
