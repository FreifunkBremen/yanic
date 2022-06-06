package prometheus_sd

import (
	"errors"

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
	TargetAddressIP     TargetAddressType = "ip"
	TargetAddressNodeID TargetAddressType = "node_id"
)

type TargetAddressFunc func(*runtime.Node) string

var TargetAddressTypeFuncs = map[TargetAddressType]TargetAddressFunc{
	TargetAddressIP: func(n *runtime.Node) string {
		if addr := n.Address; addr != nil {
			return addr.IP.String()

		}
		return ""
	},
	TargetAddressNodeID: func(n *runtime.Node) string {
		if ni := n.Nodeinfo; ni != nil {
			return ni.NodeID
		}
		return ""
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
	Targets []string               `json:"targets"`
	Labels  map[string]interface{} `json:"labels,omitempty"`
}

func toTargets(n *runtime.Node, defaultLabels map[string]interface{}, targetFunc TargetAddressFunc) *Targets {
	target := targetFunc(n)
	if target == "" {
		return nil
	}

	labels := map[string]interface{}{}
	for k, v := range defaultLabels {
		labels[k] = v
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

		// wireless - airtime
		if wifi := ni.Wireless; wifi != nil {
			labels["wifi_txpower24"] = wifi.TxPower24
			labels["wifi_channel24"] = wifi.Channel24
			labels["wifi_txpower5"] = wifi.TxPower5
			labels["wifi_channel5"] = wifi.Channel5
		}
	}
	return &Targets{
		Targets: []string{target},
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
