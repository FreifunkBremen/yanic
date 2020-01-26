package webserver

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWebserver(t *testing.T) {
	assert := assert.New(t)

	config := Config{
		Bind: ":12345",
		Webroot: "/tmp",
	}

	srv := New(config)
	assert.NotNil(srv)

	go Start(srv)

	time.Sleep(time.Millisecond * 200)

	assert.Panics(func() {
		Start(srv)
	}, "not allowed to listen twice")

	srv.Close()
}
