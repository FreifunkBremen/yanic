package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	nodes           = NewNodes()
	outputFile      string
	collectInterval time.Duration
	saveInterval    time.Duration
)

func main() {
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

	// Wait for SIGINT or SIGTERM
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigs
	log.Println("received", sig)

	for _, c := range collectors {
		c.Close()
	}
}

func check(e error) {
	if e != nil {
		log.Panic(e)
	}
}
