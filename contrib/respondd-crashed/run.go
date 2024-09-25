package main

import (
	"encoding/json"
	"net"
	"os"
	"sync"

	"github.com/bdlm/log"
	"github.com/digineo/go-ping"

	meshviewerFFRGB "github.com/FreifunkBremen/yanic/output/meshviewer-ffrgb"
)

func pingNode(pinger *ping.Pinger, node *meshviewerFFRGB.Node, addrStr string) bool {
	logNode := log.WithField("node_id", node.NodeID)

	addr, err := net.ResolveIPAddr("ip6", addrStr)
	if err != nil {
		logNode.Warnf("error parse ip address for ping: %s", err)
	}

	if addrStr[:5] == "fe80:" {
		if iface == "" {
			logNode.Debug("skip ll-addr")
			return false
		}
		addr.Zone = iface
	}
	logNode = logNode.WithField("addr", addr.String())

	_, err = pinger.PingAttempts(addr, pingTimeout, pingCount)

	logNode.WithFields(map[string]interface{}{
		"success": err == nil,
	}).Debug("pong")
	return err == nil
}

func run(pinger *ping.Pinger) {
	status := &Status{NodesCrashed: []*Node{}}
	var meshviewerjson meshviewerFFRGB.Meshviewer

	if meshviewerPATH[:4] == "http" {
		if err := JSONRequest(meshviewerPATH, &meshviewerjson); err != nil {
			status.Error = err.Error()
			log.Errorf("error during fetch meshviewer.json: %s", err)
		}
	} else {
		meshviewerFile, err := os.Open(meshviewerPATH)
		if err != nil {
			status.Error = err.Error()
			log.Errorf("error during fetch meshviewer.json: %s", err)
		} else if err := json.NewDecoder(meshviewerFile).Decode(&meshviewerjson); err != nil {
			status.Error = err.Error()
			log.Errorf("error during decode meshviewer.json: %s", err)
		}
	}

	log.Debug("fetched meshviewer.json")

	wg := sync.WaitGroup{}
	wg.Add(len(meshviewerjson.Nodes))

	offline := 0
	for _, node := range meshviewerjson.Nodes {
		go func(node *meshviewerFFRGB.Node) {
			defer wg.Done()
			if node.IsOnline {
				return
			}
			logNode := log.WithField("node", node.NodeID)
			wgNode := sync.WaitGroup{}
			wgNode.Add(len(node.Addresses))
			offline += 1
			notReachable := true
			for _, addr := range node.Addresses {
				go func(node *meshviewerFFRGB.Node, addr string) {
					if ok := pingNode(pinger, node, addr); ok {
						notReachable = false
					}
					wgNode.Done()
				}(node, addr)
			}
			wgNode.Wait()
			if !notReachable {
				logNode.Info("add to crashed list")
				status.AddNode(node)
			}
		}(node)
	}

	wg.Wait()

	status.Lock()
	status.NodesCount = len(meshviewerjson.Nodes)
	status.NodesOfflineCount = offline
	status.Unlock()

	tmpFile := statusPath + ".tmp"
	statusFile, err := os.OpenFile(tmpFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Warnf("unable to open status file: %s", err)
	}
	defer statusFile.Close()

	if err := json.NewEncoder(statusFile).Encode(status); err != nil {
		log.Warnf("unable to write status json: %s", err)
	}
	if err := os.Rename(tmpFile, statusPath); err != nil {
		log.Warnf("unable to move status file: %s", err)
	}

	log.WithFields(map[string]interface{}{
		"count_meshviewer": status.NodesCount,
		"count_offline":    status.NodesOfflineCount,
		"count_status":     len(status.NodesCrashed),
	}).Info("test complete")
}
