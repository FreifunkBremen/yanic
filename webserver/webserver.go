package webserver

import (
	"net/http"

	"github.com/NYTimes/gziphandler"
	"github.com/bdlm/log"
)

// New creates a new webserver and starts it
func New(config Config) *http.Server {
	return &http.Server{
		Addr:    config.Bind,
		Handler: gziphandler.GzipHandler(http.FileServer(http.Dir(config.Webroot))),
	}
}

func Start(srv *http.Server) {
	// service connections
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Panicf("webserver crashed: %s", err)
	}
}
