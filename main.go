package main

import (
	"flag"
	"log"
	"time"
	"net/http"
)
var (
	nodeserver	= NewNodeServer("/nodes")
	nodes           = NewNodes()
	outputFile      string
	collectInterval time.Duration
	saveInterval    time.Duration
)
func main(){
	var collectSeconds, saveSeconds int

	flag.StringVar(&outputFile, "output", "nodes.json", "path output file")
	flag.IntVar(&collectSeconds, "collectInterval", 15, "interval for data collections")
	flag.IntVar(&saveSeconds, "saveInterval", 5, "interval for data saving")
	flag.Parse()

	collectInterval = time.Second * time.Duration(collectSeconds)
	saveInterval = time.Second * time.Duration(saveSeconds)

	collectors := []*Collector{
		NewCollector("statistics"),
		NewCollector("nodeinfo"),
		NewCollector("neighbours"),
	}	

	go nodeserver.Listen()
	
	// static files
	http.Handle("/", http.FileServer(http.Dir("webroot")))

	log.Fatal(http.ListenAndServe(":8080", nil))
	for _, c := range collectors {
		c.Close()
	}
}
