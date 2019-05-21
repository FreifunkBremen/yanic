package geojson

import (
	"errors"

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
	output.RegisterAdapter("geojson", Register)
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

	runtime.SaveJSON(transform(nodes), o.path)
}
