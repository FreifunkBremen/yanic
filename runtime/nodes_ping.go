package runtime

import (
	"github.com/bdlm/log"
	"github.com/sparrc/go-ping"
)

func (nodes *Nodes) ping(node *Node) bool {
	logNode := log.WithField("node_id", "unknown")
	if node.Nodeinfo != nil {
		logNode = logNode.WithField("node_id", node.Nodeinfo.NodeID)
	}
	var addr string
	if node.Address != nil {
		addr = node.Address.IP.String()
		if node.Address.IP.IsLinkLocalUnicast() {
			addr += "%" + node.Address.Zone
		}
	} else {
		logNode.Debug("error no address found")
		if node.Nodeinfo != nil {
			for _, addrMaybe := range node.Nodeinfo.Network.Addresses {
				if len(addrMaybe) >= 5 && addrMaybe[:5] != "fe80:" {
					addr = addrMaybe
				}
			}
		}
	}

	logAddr := logNode.WithField("addr", addr)

	pinger, err := ping.NewPinger(addr)
	if err != nil {
		logAddr.Debugf("error during ping: %s", err)
		return false
	}
	pinger.SetPrivileged(true)
	pinger.Count = nodes.config.PingCount
	pinger.Timeout = nodes.config.PingTimeout.Duration
	pinger.Interval =  pinger.Timeout / pinger.Count
	pinger.Run() // blocks until finished
	stats := pinger.Statistics()
	logAddr.WithFields(map[string]interface{}{
		"pkg_lost": stats.PacketLoss,
	}).Debug("pong")
	return stats.PacketLoss < 100
}
