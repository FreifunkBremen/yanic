package meshviewer

import (
	"log"

	"github.com/FreifunkBremen/yanic/output"
	"github.com/FreifunkBremen/yanic/runtime"
)

type Output struct {
	output.Output
	config  Config
	nodes   *runtime.Nodes
	builder nodeBuilder
	filter  filter
}

type Config map[string]interface{}

func (c Config) Enable() bool {
	return c["enable"].(bool)
}

func (c Config) Version() int64 {
	return c["version"].(int64)
}
func (c Config) NodesPath() string {
	if c["nodes_path"] == nil {
		log.Panic("in configuration of [[nodes.output.meshviewer]] was no nodes_path defined", c)
	}
	return c["nodes_path"].(string)
}
func (c Config) GraphPath() string {
	return c["graph_path"].(string)
}

func (c Config) FilterOption() filterConfig {
	if v, ok := c["filter"]; ok {
		var filterMap filterConfig
		filterMap = v.(map[string]interface{})
		return filterMap
	}
	return nil
}

type nodeBuilder func(filter, *runtime.Nodes) interface{}

var nodeFormats = map[int64]nodeBuilder{
	1: BuildNodesV1,
	2: BuildNodesV2,
}

func init() {
	output.RegisterAdapter("meshviewer", Register)
}

func Register(nodes *runtime.Nodes, configuration interface{}) (output.Output, error) {
	var config Config
	config = configuration.(map[string]interface{})
	if !config.Enable() {
		return nil, nil
	}

	builder := nodeFormats[config.Version()]
	if builder == nil {
		log.Panicf("invalid nodes version: %d", config.Version())
	}

	return &Output{
		nodes:   nodes,
		config:  config,
		builder: builder,
		filter:  createFilter(config.FilterOption()),
	}, nil
}

func (o *Output) Save() {
	o.nodes.RLock()
	defer o.nodes.RUnlock()

	if path := o.config.NodesPath(); path != "" {
		runtime.SaveJSON(o.builder(o.filter, o.nodes), path)
	}

	if path := o.config.GraphPath(); path != "" {
		runtime.SaveJSON(BuildGraph(o.nodes), path)
	}
}
