package all

import (
	"github.com/FreifunkBremen/yanic/output"
	"github.com/FreifunkBremen/yanic/runtime"
)

type Output struct {
	output.Output
	nodes *runtime.Nodes
	list  []output.Output
}

func Register(nodes *runtime.Nodes, configuration interface{}) (output.Output, error) {
	var list []output.Output
	allOutputs := configuration.(map[string][]interface{})
	for outputType, outputRegister := range output.Adapters {
		outputConfigs := allOutputs[outputType]
		for _, config := range outputConfigs {
			output, err := outputRegister(nodes, config)
			if err != nil {
				return nil, err
			}
			if output == nil {
				continue
			}
			list = append(list, output)
		}
	}
	return &Output{list: list, nodes: nodes}, nil
}

func (o *Output) Save() {
	for _, item := range o.list {
		item.Save()
	}
}
