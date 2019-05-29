package prometheus

import (
	"errors"
	"os"

	"github.com/bdlm/log"

	"github.com/FreifunkBremen/yanic/output"
	"github.com/FreifunkBremen/yanic/runtime"
)

type Output struct {
	output.Output
	path string
}

type Config map[string]interface{}

func (c Config) Path() string {
	if path, ok := c["path"]; ok {
		return path.(string)
	}
	return ""
}

func init() {
	output.RegisterAdapter("prometheus", Register)
}

func Register(configuration map[string]interface{}) (output.Output, error) {
	var config Config
	config = configuration

	if path := config.Path(); path != "" {
		return &Output{
			path: path,
		}, nil
	}
	return nil, errors.New("no path given")

}

func (o *Output) Save(nodes *runtime.Nodes) {
	nodes.RLock()
	defer nodes.RUnlock()

	tmpFile := o.path + ".tmp"

	f, err := os.OpenFile(tmpFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Panic(err)
	}

	for _, node := range nodes.List {
		metrics := MetricsFromNode(nodes, node)
		for _, m := range metrics {
			str, err := m.String()
			if err == nil {
				f.WriteString(str + "\n")
			} else {
				logger := log.WithField("database", "prometheus")
				if nodeinfo := node.Nodeinfo; nodeinfo != nil {
					logger = logger.WithField("node_id", nodeinfo.NodeID)
				}
				logger.Warnf("not able to get metrics from node: %s", err)
			}
		}
	}
	f.Close()

	if err := os.Rename(tmpFile, o.path); err != nil {
		log.Panic(err)
	}
}
