package webserver

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWebserver(t *testing.T) {
	assert := assert.New(t)

	srv := New(":12345", "/tmp")
	assert.NotNil(srv)

	go Start(srv)

	time.Sleep(time.Millisecond * 200)

	assert.Panics(func() {
		Start(srv)
	}, "not allowed to listen twice")

	if err := srv.Close(); err != nil {
		fmt.Println("Error when closing:", err)
	}
}
