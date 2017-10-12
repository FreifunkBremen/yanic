package webserver

import (
	"net/http"

	"github.com/NYTimes/gziphandler"
)

// New creates a new webserver and starts it
func New(bindAddr, webroot string) *http.Server {
	return &http.Server{
		Addr:    bindAddr,
		Handler: gziphandler.GzipHandler(http.FileServer(http.Dir(webroot))),
	}
}

func Start(srv *http.Server) {
	// service connections
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		panic(err)
	}
}
