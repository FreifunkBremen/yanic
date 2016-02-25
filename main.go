package main

import (
	"encoding/json"
	"flag"
	"log"
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
	wsserverForNodes = websocketserver.NewServer("/nodes")
	responedDaemon   *responed.Daemon
	nodes            = models.NewNodes()
	outputFile       string
	collectInterval  time.Duration
	saveInterval     time.Duration
)

func main() {
	var collectSeconds, saveSeconds int

	flag.StringVar(&outputFile, "output", "webroot/nodes.json", "path output file")
	flag.IntVar(&collectSeconds, "collectInterval", 15, "interval for data collections")
	flag.IntVar(&saveSeconds, "saveInterval", 5, "interval for data saving")
	flag.Parse()

	collectInterval = time.Second * time.Duration(collectSeconds)
	saveInterval = time.Second * time.Duration(saveSeconds)

	go wsserverForNodes.Listen()
	go nodes.Saver(outputFile, saveInterval)
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
			log.Printf("Start sendAll after recieve %s", coll.CollectType)
			log.Printf("Field: %V", field)
			log.Printf("Result: %V", result)
			wsserverForNodes.SendAll(node)
			log.Print("End")
		}

		field.Set(reflect.ValueOf(result))
	})
	go responedDaemon.ListenAndSend(collectInterval)

	http.Handle("/", http.FileServer(http.Dir("webroot")))
	//TODO bad
	log.Fatal(http.ListenAndServe(":8080", nil))

	// Wait for End
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigs
	log.Println("received", sig)

	// Close everything at the end
	wsserverForNodes.Close()
	responedDaemon.Close()
}
