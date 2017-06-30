package template

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os"
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

func (c Config) TemplatePath() string {
	return c["template_path"].(string)
}
func (c Config) ResultPath() string {
	return c["result_path"].(string)
}

func init() {
	output.RegisterAdapter("template", Register)
}

func Register(nodes *runtime.Nodes, configuration interface{}) (output.Output, error) {
	var config Config
	config = configuration.(map[string]interface{})
	if !config.Enable() {
		return nil, nil
	}
	t := goTemplate.New("some")
	t = t.Funcs(goTemplate.FuncMap{"json": func(v interface{}) string {
		a, _ := json.Marshal(v)
		return string(a)
	}})
	buf := bytes.NewBuffer(nil)
	f, err := os.Open(config.TemplatePath()) // Error handling elided for brevity.
	if err != nil {
		log.Panic(err)
	}
	io.Copy(buf, f) // Error handling elided for brevity.
	f.Close()

	s := string(buf.Bytes())
	t.Parse(s)
	return &Output{
		config:   config,
		nodes:    nodes,
		template: t,
	}, nil
}

func (o *Output) Save() {
	stats := runtime.NewGlobalStats(o.nodes)
	if stats == nil {
		log.Panic("update of [output.template] not possible invalid data for the template generated")
	}
	tmpFile := o.config.ResultPath() + ".tmp"
	f, err := os.OpenFile(tmpFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Panic(err)
	}
	o.template.Execute(f, map[string]interface{}{"GlobalStatistic": stats})
	if err != nil {
		log.Panic(err)
	}
	f.Close()
	if err := os.Rename(tmpFile, o.config.ResultPath()); err != nil {
		log.Panic(err)
	}
}
