package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	collector := NewCollector()
	defer collector.Close()
	collector.send("[2a06:8782:ffbb:1337:c24a:ff:fe2c:c7ac]:1001")
	collector.send("[2001:bf7:540:0:32b5:c2ff:fe6e:99d5]:1001")

	// Wait for SIGINT or SIGTERM
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigs
	log.Println("received", sig)

}
