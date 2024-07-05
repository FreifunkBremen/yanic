package webserver

import (
	"net/http"

	"github.com/FreifunkBremen/yanic/runtime"

	"github.com/NYTimes/gziphandler"
	"github.com/bdlm/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// New creates a new webserver and starts it
func New(config Config, nodes *runtime.Nodes) *http.Server {
	mux := http.NewServeMux()
	if config.Prometheus != nil && config.Prometheus.Enable {
		config.Prometheus.Init(nodes)
		prometheus.MustRegister(config.Prometheus)
		mux.Handle("/metrics", promhttp.Handler())
	}
	if config.Webroot != "" {
		mux.Handle("/", gziphandler.GzipHandler(http.FileServer(http.Dir(config.Webroot))))
	}
	return &http.Server{
		Addr:    config.Bind,
		Handler: mux,
	}
}

func Start(srv *http.Server) {
	// service connections
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Panicf("webserver crashed: %s", err)
	}
}
