package webserver

import (
	"net/http"

	"github.com/NYTimes/gziphandler"
)

// New creates a new webserver and starts it
func New(bindAddr, webroot string) *http.Server {
	srv := &http.Server{
		Addr:    bindAddr,
		Handler: gziphandler.GzipHandler(http.FileServer(http.Dir(webroot))),
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil {
			panic(err)
		}
	}()

	return srv
}
