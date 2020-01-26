package webserver

import (
	"io"
	"net/http"

	"github.com/FreifunkBremen/yanic/lib/duration"

	"github.com/FreifunkBremen/yanic/respond"
	"github.com/FreifunkBremen/yanic/runtime"
)

type PrometheusConfig struct {
	Enable   bool              `toml:"enable"`
	Wait     duration.Duration `toml:"wait"`
	Outdated duration.Duration `toml:"outdated"`
}

func CreatePrometheusExporter(config PrometheusConfig, srv *http.Server, coll *respond.Collector, nodes *runtime.Nodes) {
	mux := http.NewServeMux()
	mux.HandleFunc("/metric", func(res http.ResponseWriter, req *http.Request) {
		io.WriteString(res, "Hello from a HandleFunc #2!\n")
	})
	if srv.Handler != nil {
		mux.Handle("/", srv.Handler)
	}
	srv.Handler = mux
}
