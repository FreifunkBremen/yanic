package influxdb

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	assert := assert.New(t)

	conn, err := Connect(map[string]interface{}{
		"address":        "",
		"token":          "",
		"bucket_default": "all",
	})
	assert.Nil(conn)
	assert.Error(err)

	conn, err = Connect(map[string]interface{}{
		"address":        "http://localhost",
		"token":          "",
		"bucket_default": "all",
	})
	assert.Nil(conn)
	assert.Error(err)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	conn, err = Connect(map[string]interface{}{
		"address":        srv.URL,
		"token":          "atoken",
		"bucket_default": "all",
	})

	assert.NotNil(conn)
	assert.NoError(err)
}
