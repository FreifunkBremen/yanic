package runtime

import (
	"net"

	"github.com/bdlm/log"
)

func (nodes *Nodes) ping(node *Node) bool {
	logNode := log.WithField("node_id", "unknown")
	if node.Nodeinfo != nil {
		logNode = logNode.WithField("node_id", node.Nodeinfo.NodeID)
	}
	var addr *net.IPAddr
	if node.Address != nil {
		addr = &net.IPAddr{IP:node.Address.IP, Zone: node.Address.Zone}
	} else {
		logNode.Debug("error no address found")
		if node.Nodeinfo != nil {
			for _, addrMaybeString := range node.Nodeinfo.Network.Addresses {
				if len(addrMaybeString) >= 5 && addrMaybeString[:5] != "fe80:" {
					addrMaybe, err := net.ResolveIPAddr("ip6", addrMaybeString)
					if err == nil {
						addr = addrMaybe
					}
				}
			}
		}
	}

	logAddr := logNode.WithField("addr", addr.String())

	_, err := nodes.pinger.PingAttempts(addr, nodes.config.PingTimeout.Duration, nodes.config.PingCount)

	logAddr.WithFields(map[string]interface{}{
		"success": err == nil,
	}).Debug("pong")
	return err == nil
}
