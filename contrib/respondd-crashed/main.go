package main

import (
	"flag"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/bdlm/log"
	stdLogger "github.com/bdlm/std/logger"
	"github.com/digineo/go-ping"
)

var (
	timestamps bool
	loglevel   uint

	runEvery time.Duration

	iface string

	pingCount   int
	pingTimeout time.Duration

	meshviewerPATH string
	statusPath     string
)

func main() {
	flag.BoolVar(&timestamps, "timestamps", false, "Enables timestamps for log output")
	flag.UintVar(&loglevel, "loglevel", 40, "Show log message starting at level")

	flag.DurationVar(&runEvery, "run-every", time.Duration(time.Minute), "repeat check every")

	flag.StringVar(&iface, "ll-iface", "", "interface to ping linklocal-address")

	flag.IntVar(&pingCount, "ping-count", 3, "count of pings")
	flag.DurationVar(&pingTimeout, "ping-timeout", time.Duration(time.Second*5), "timeout to wait for response")

	flag.StringVar(&statusPath, "status-path", "respondd-crashed.json", "path to store status")
	flag.StringVar(&meshviewerPATH, "meshviewer-path", "meshviewer.json", "path to meshviewer.json from yanic")

	flag.Parse()

	log.AddHook(&Hook{})
	log.SetLevel(stdLogger.Level(loglevel))
	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp: timestamps,
	})

	pinger, err := ping.New("", "::")
	if err != nil {
		log.Panicf("not able to bind pinger: %s", err)
	}

	timer := time.NewTimer(runEvery)

	stop := false

	wg := sync.WaitGroup{}

	log.Info("start tester")

	func() {
		wg.Add(1)
		for !stop {
			select {
			case <-timer.C:
				run(pinger)
				timer.Reset(runEvery)
			}
		}
		timer.Stop()
		wg.Done()
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigs
	stop = true
	wg.Wait()
	log.Infof("stopped: %s", sig)

}
