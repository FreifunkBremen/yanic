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
	var config Config
	config = configuration

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

func (o *Output) Save(nodes *runtime.Nodes) {
	nodes.RLock()
	defer nodes.RUnlock()

	targets := &Targets{
		Targets: []string{},
		Labels:  o.labels,
	}
	if o.targetType == TargetAddressNodeID {
		for _, n := range nodes.List {
			if ni := n.Nodeinfo; ni != nil {
				targets.Targets = append(targets.Targets, ni.NodeID)
			}
		}
	} else {
		for _, n := range nodes.List {
			if addr := n.Address; addr != nil {
				targets.Targets = append(targets.Targets, addr.IP.String())
			}
		}
	}

	runtime.SaveJSON([]interface{}{targets}, o.path)
}
