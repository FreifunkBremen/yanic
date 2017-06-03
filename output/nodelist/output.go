package nodelist

import (
	goTemplate "text/template"

	"github.com/FreifunkBremen/yanic/output"
	"github.com/FreifunkBremen/yanic/runtime"
)

type Output struct {
	output.Output
	config   Config
	nodes    *runtime.Nodes
	template *goTemplate.Template
}

type Config map[string]interface{}

func (c Config) Enable() bool {
	return c["enable"].(bool)
}

func (c Config) Path() string {
	return c["path"].(string)
}

func init() {
	output.RegisterAdapter("nodelist", Register)
}

func Register(nodes *runtime.Nodes, configuration interface{}) (output.Output, error) {
	var config Config
	config = configuration.(map[string]interface{})
	if !config.Enable() {
		return nil, nil
	}

	return &Output{
		config: config,
		nodes:  nodes,
	}, nil
}

func (o *Output) Save() {
	o.nodes.RLock()
	defer o.nodes.RUnlock()

	if path := o.config.Path(); path != "" {
		runtime.SaveJSON(transform(o.nodes), path)
	}
}
