package meshviewer

import (
	"fmt"

	"github.com/bdlm/log"

	"github.com/FreifunkBremen/yanic/output"
	"github.com/FreifunkBremen/yanic/runtime"
)

type Output struct {
	output.Output
	config  Config
	builder nodeBuilder
}

type Config map[string]interface{}

func (c Config) Version() int64 {
	if v := c["version"]; v != nil {
		return v.(int64)
	}
	return -1
}
func (c Config) NodesPath() string {
	if c["nodes_path"] == nil {
		log.Panic("in configuration of [[nodes.output.meshviewer]] was no nodes_path defined:\n", c)
	}
	return c["nodes_path"].(string)
}
func (c Config) GraphPath() string {
	return c["graph_path"].(string)
}

type nodeBuilder func(*runtime.Nodes) interface{}

var nodeFormats = map[int64]nodeBuilder{
	1: BuildNodesV1,
	2: BuildNodesV2,
}

func init() {
	output.RegisterAdapter("meshviewer", Register)
}

func Register(configuration map[string]interface{}) (output.Output, error) {
	var config Config
	config = configuration

	builder := nodeFormats[config.Version()]
	if builder == nil {
		return nil, fmt.Errorf("invalid nodes version: %d", config.Version())
	}

	return &Output{
		config:  config,
		builder: builder,
	}, nil
}

func (o *Output) Save(nodes *runtime.Nodes) {
	nodes.RLock()
	defer nodes.RUnlock()

	if path := o.config.NodesPath(); path != "" {
		runtime.SaveJSON(o.builder(nodes), path)
	}

	if path := o.config.GraphPath(); path != "" {
		runtime.SaveJSON(BuildGraph(nodes), path)
	}
}
