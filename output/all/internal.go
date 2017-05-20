package all

import (
	"log"

	"github.com/FreifunkBremen/yanic/output"
	"github.com/FreifunkBremen/yanic/runtime"
)

type Output struct {
	output.Output
	list   map[int]output.Output
	filter map[int]filterConfig
}

func Register(configuration map[string]interface{}) (output.Output, error) {
	list := make(map[int]output.Output)
	filter := make(map[int]filterConfig)
	i := 1
	allOutputs := configuration
	for outputType, outputRegister := range output.Adapters {
		configForOutput := allOutputs[outputType]
		if configForOutput == nil {
			log.Printf("the output type '%s' has no configuration\n", outputType)
			continue
		}
		outputConfigs, ok := configForOutput.([]map[string]interface{})
		if !ok {
			log.Panicf("the output type '%s' has the wrong format\n", outputType)
		}
		for _, config := range outputConfigs {
			if c, ok := config["enable"].(bool); ok && !c {
				continue
			}
			output, err := outputRegister(config)
			if err != nil {
				return nil, err
			}
			if output == nil {
				continue
			}
			list[i] = output
			if c := config["filter"]; c != nil {
				filter[i] = config["filter"].(map[string]interface{})
			}
			i++
		}
	}
	return &Output{list: list, filter: filter}, nil
}

func (o *Output) Save(nodes *runtime.Nodes) {
	for i, item := range o.list {
		var filteredNodes *runtime.Nodes
		if config := o.filter[i]; config != nil {
			filteredNodes = config.filtering(nodes)
		} else {
			filteredNodes = filterConfig{}.filtering(nodes)
		}

		item.Save(filteredNodes)
	}
}
