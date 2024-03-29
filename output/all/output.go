package all

import (
	"fmt"

	"github.com/bdlm/log"

	"github.com/FreifunkBremen/yanic/output"
	"github.com/FreifunkBremen/yanic/output/filter"
	"github.com/FreifunkBremen/yanic/runtime"
)

type Output struct {
	output.Output
	list         map[int]output.Output
	outputFilter map[int]filter.Set
}

func Register(configuration map[string]interface{}) (output.Output, error) {
	list := make(map[int]output.Output)
	outputFilter := make(map[int]filter.Set)
	i := 1
	allOutputs := configuration
	for outputType, outputRegister := range output.Adapters {
		configForOutput := allOutputs[outputType]
		if configForOutput == nil {
			log.WithField("output", outputType).Infof("no configuration found")
			continue
		}
		outputConfigs, ok := configForOutput.([]map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("the output type '%s' has the wrong format: read format %T should be an []map[string]interface{}", outputType, configForOutput)
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
			var errs []error
			var filterSet filter.Set
			if c := config["filter"]; c != nil {
				if filterConf, ok := c.(map[string]interface{}); ok {
					filterSet, errs = filter.New(filterConf)
				}
				if len(errs) > 0 {
					return nil, fmt.Errorf("filter configuration errors: %v", errs)
				}
				outputFilter[i] = filterSet
			}
			list[i] = output
			i++
		}
	}
	return &Output{list: list, outputFilter: outputFilter}, nil
}

func (o *Output) Save(nodes *runtime.Nodes) {
	for i, item := range o.list {
		item.Save(o.outputFilter[i].Apply(nodes))
	}
}
