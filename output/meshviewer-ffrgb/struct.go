package meshviewerFFRGB

import (
	"time"

	"github.com/FreifunkBremen/yanic/lib/jsontime"
	"github.com/FreifunkBremen/yanic/runtime"
)

type Meshviewer struct {
	Timestamp jsontime.Time `json:"timestamp"`
	Nodes     []*Node       `json:"nodes"`
	Links     []*Link       `json:"links"`
}

type Node struct {
	Firstseen      jsontime.Time `json:"firstseen"`
	Lastseen       jsontime.Time `json:"lastseen"`
	IsOnline       bool          `json:"is_online"`
	IsGateway      bool          `json:"is_gateway"`
	Clients        uint32        `json:"clients"`
	ClientsWifi24  uint32        `json:"clients_wifi24"`
	ClientsWifi5   uint32        `json:"clients_wifi5"`
	ClientsOthers  uint32        `json:"clients_other"`
	RootFSUsage    float64       `json:"rootfs_usage"`
	LoadAverage    float64       `json:"loadavg"`
	MemoryUsage    *float64      `json:"memory_usage,omitempty"`
	Uptime         jsontime.Time `json:"uptime,omitempty"`
	GatewayNexthop string        `json:"gateway_nexthop,omitempty"`
	GatewayIPv4    string        `json:"gateway,omitempty"`
	GatewayIPv6    string        `json:"gateway6,omitempty"`
	NodeID         string        `json:"node_id"`
	MAC            string        `json:"mac"`
	Addresses      []string      `json:"addresses"`
	SiteCode       string        `json:"-"`
	DomainCode     string        `json:"domain"`
	Hostname       string        `json:"hostname"`
	Owner          string        `json:"owner,omitempty"`
	Location       *Location     `json:"location,omitempty"`
	Firmware       Firmware      `json:"firmware,omitempty"`
	Autoupdater    Autoupdater   `json:"autoupdater"`
	Nproc          int           `json:"nproc"`
	Model          string        `json:"model,omitempty"`
}

// Firmware out of software
type Firmware struct {
	Base    string `json:"base,omitempty"`
	Release string `json:"release,omitempty"`
}

// Autoupdater
type Autoupdater struct {
	Enabled bool   `json:"enabled"`
	Branch  string `json:"branch,omitempty"`
}

// Location struct
type Location struct {
	Longitude float64 `json:"longitude,omitempty"`
	Latitude  float64 `json:"latitude,omitempty"`
}

// Link
type Link struct {
	Type          string  `json:"type"`
	Source        string  `json:"source"`
	Target        string  `json:"target"`
	SourceTQ      float32 `json:"source_tq"`
	TargetTQ      float32 `json:"target_tq"`
	SourceAddress string  `json:"source_addr"`
	TargetAddress string  `json:"target_addr"`
}

func NewNode(nodes *runtime.Nodes, n *runtime.Node) *Node {
	node := &Node{
		Firstseen: n.Firstseen,
		Lastseen:  n.Lastseen,
		IsOnline:  n.Online,
		IsGateway: n.IsGateway(),
	}

	if nodeinfo := n.Nodeinfo; nodeinfo != nil {
		node.NodeID = nodeinfo.NodeID
		node.MAC = nodeinfo.Network.Mac
		node.Addresses = nodeinfo.Network.Addresses
		node.SiteCode = nodeinfo.System.SiteCode
		node.DomainCode = nodeinfo.System.DomainCode
		node.Hostname = nodeinfo.Hostname
		if owner := nodeinfo.Owner; owner != nil {
			node.Owner = owner.Contact
		}
		if location := nodeinfo.Location; location != nil {
			node.Location = &Location{
				Longitude: location.Longitude,
				Latitude:  location.Latitude,
			}
		}
		node.Firmware = nodeinfo.Software.Firmware
		node.Autoupdater = Autoupdater{
			Enabled: nodeinfo.Software.Autoupdater.Enabled,
			Branch:  nodeinfo.Software.Autoupdater.Branch,
		}
		node.Nproc = nodeinfo.Hardware.Nproc
		node.Model = nodeinfo.Hardware.Model
	}
	if statistic := n.Statistics; statistic != nil {
		if n.Online {
			node.Clients = statistic.Clients.Total
			node.ClientsWifi24 = statistic.Clients.Wifi24
			node.ClientsWifi5 = statistic.Clients.Wifi5

			clientsWifi := node.ClientsWifi24 + node.ClientsWifi5
			if node.Clients == 0 {
				node.Clients = clientsWifi
			} else if node.Clients >= clientsWifi {
				node.ClientsOthers = node.Clients - clientsWifi
			}
		}

		node.RootFSUsage = statistic.RootFsUsage
		node.LoadAverage = statistic.LoadAverage

		/* The Meshviewer could not handle absolute memory output
		 * calc the used memory as a float which 100% equal 1.0
		 * calc is coppied from node statuspage (look discussion:
		 * https://github.com/FreifunkBremen/yanic/issues/35 and
		 * https://github.com/freifunk-gluon/gluon/pull/1517)
		 */
		if statistic.Memory.Total > 0 {
			usage := 0.0
			if statistic.Memory.Available > 0 {
				usage = 1 - float64(statistic.Memory.Available)/float64(statistic.Memory.Total)
			} else {
				usage = 1 - (float64(statistic.Memory.Free)+float64(statistic.Memory.Buffers)+float64(statistic.Memory.Cached))/float64(statistic.Memory.Total)
			}
			node.MemoryUsage = &usage
		}

		node.Uptime = jsontime.Now().Add(time.Duration(statistic.Uptime) * -time.Second)
		node.GatewayNexthop = nodes.GetNodeIDbyAddress(statistic.GatewayNexthop)
		if node.GatewayNexthop == "" {
			node.GatewayNexthop = statistic.GatewayNexthop
		}
		node.GatewayIPv4 = nodes.GetNodeIDbyAddress(statistic.GatewayIPv4)
		if node.GatewayIPv4 == "" {
			node.GatewayIPv4 = statistic.GatewayIPv4
		}
		node.GatewayIPv6 = nodes.GetNodeIDbyAddress(statistic.GatewayIPv6)
		if node.GatewayIPv6 == "" {
			node.GatewayIPv6 = statistic.GatewayIPv6
		}
	}

	// fix site to domain - if empty
	if node.DomainCode == "" {
		node.DomainCode = node.SiteCode
	}

	return node
}
