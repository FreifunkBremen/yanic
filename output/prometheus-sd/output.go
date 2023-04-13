package prometheus_sd

import (
	"errors"
	"fmt"
	"net"
	"strconv"

	"github.com/FreifunkBremen/yanic/output"
	"github.com/FreifunkBremen/yanic/runtime"
)

type Output struct {
	output.Output
	path       string
	targetType TargetAddressType
	labels     map[string]interface{}
}

type Config map[string]interface{}

func (c Config) Path() string {
	if path, ok := c["path"]; ok {
		return path.(string)
	}
	return ""
}

type TargetAddressType string

const (
	TargetAddressNodeID    TargetAddressType = "node_id"
	TargetAddressIP        TargetAddressType = "ip"
	TargetAddressIPPublish TargetAddressType = "ip-publish"
)

type TargetAddressFunc func(*runtime.Node) []string

var TargetAddressTypeFuncs = map[TargetAddressType]TargetAddressFunc{
	TargetAddressNodeID: func(n *runtime.Node) []string {
		if ni := n.Nodeinfo; ni != nil {
			return []string{ni.NodeID}
		}
		return []string{}
	},
	TargetAddressIP: func(n *runtime.Node) []string {
		if addr := n.Address; addr != nil {
			return []string{addr.IP.String()}

		}
		return []string{}
	},
	TargetAddressIPPublish: func(n *runtime.Node) []string {
		addresses := []string{}
		if nodeinfo := n.Nodeinfo; nodeinfo != nil {
			for _, addr := range nodeinfo.Network.Addresses {
				if net.ParseIP(addr).IsGlobalUnicast() {
					addresses = append(addresses, addr)
				}
			}

		}
		return addresses
	},
}

func (c Config) TargetAddress() TargetAddressType {
	if v, ok := c["target_address"]; ok {
		return TargetAddressType(v.(string))
	}
	return TargetAddressIP
}

func (c Config) Labels() map[string]interface{} {
	if v, ok := c["labels"]; ok {
		return v.(map[string]interface{})
	}
	return nil
}

func init() {
	output.RegisterAdapter("prometheus-sd", Register)
}

func Register(configuration map[string]interface{}) (output.Output, error) {
	config := Config(configuration)

	if path := config.Path(); path != "" {
		return &Output{
			path:       path,
			targetType: config.TargetAddress(),
			labels:     config.Labels(),
		}, nil
	}
	return nil, errors.New("no path given")

}

type Targets struct {
	Targets []string          `json:"targets"`
	Labels  map[string]string `json:"labels,omitempty"`
}

func toTargets(n *runtime.Node, defaultLabels map[string]interface{}, targetFunc TargetAddressFunc) *Targets {
	target := targetFunc(n)
	if len(target) <= 0 {
		return nil
	}

	labels := map[string]string{}
	for k, v := range defaultLabels {
		vS := v.(string)
		labels[k] = vS
	}
	if ni := n.Nodeinfo; ni != nil {
		labels["node_id"] = ni.NodeID
		labels["hostname"] = ni.Hostname
		// model
		if model := ni.Hardware.Model; model != "" {
			labels["model"] = model
		}
		// system
		if siteCode := ni.System.SiteCode; siteCode != "" {
			labels["site_code"] = siteCode
		}
		if domainCode := ni.System.DomainCode; domainCode != "" {
			labels["domain_code"] = domainCode
		}
		if primaryDomainCode := ni.System.PrimaryDomainCode; primaryDomainCode != "" {
			labels["primary_domain_code"] = primaryDomainCode
		}

		// owner
		if owner := ni.Owner; owner != nil {
			labels["owner"] = owner.Contact
		}

		// location
		if location := ni.Location; location != nil {
			labels["latitude"] = fmt.Sprintf("%v", location.Latitude)
			labels["longitude"] = fmt.Sprintf("%v", location.Longitude)
		}

		// wireless - airtime
		if wifi := ni.Wireless; wifi != nil {
			labels["wifi_txpower24"] = strconv.Itoa(int(wifi.TxPower24))
			labels["wifi_channel24"] = strconv.Itoa(int(wifi.Channel24))
			labels["wifi_txpower5"] = strconv.Itoa(int(wifi.TxPower5))
			labels["wifi_channel5"] = strconv.Itoa(int(wifi.Channel5))
		}
	}
	return &Targets{
		Targets: target,
		Labels:  labels,
	}
}

func (o *Output) Save(nodes *runtime.Nodes) {
	nodes.RLock()
	defer nodes.RUnlock()

	targetFunc, ok := TargetAddressTypeFuncs[o.targetType]
	if !ok {
		return
	}
	targets := []*Targets{}
	for _, n := range nodes.List {
		if t := toTargets(n, o.labels, targetFunc); t != nil {
			targets = append(targets, t)
		}
	}

	runtime.SaveJSON(targets, o.path)
}
