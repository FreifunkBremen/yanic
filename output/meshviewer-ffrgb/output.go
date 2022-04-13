package meshviewerFFRGB

import (
	"errors"

	"yanic/output"
	"yanic/runtime"
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
	output.RegisterAdapter("meshviewer-ffrgb", Register)
}

func Register(configuration map[string]interface{}) (output.Output, error) {
	config := Config(configuration)

	if path := config.Path(); path != "" {
		return &Output{
			path: path,
		}, nil
	}
	return nil, errors.New("no path given")

}

func (o *Output) Save(nodes *runtime.Nodes) {
	runtime.SaveJSON(transform(nodes), o.path)
}
