package logging

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/FreifunkBremen/yanic/data"
	"github.com/FreifunkBremen/yanic/runtime"
	"github.com/stretchr/testify/assert"
)

func TestStart(t *testing.T) {
	assert := assert.New(t)

	conn, err := Connect(map[string]interface{}{
		"path": "/dev/notexists/file",
	})
	assert.Nil(conn)
	assert.Error(err)

	path := "/tmp/testlogfile"

	conn, err = Connect(map[string]interface{}{
		"path": path,
	})
	assert.NoError(err)

	dat, _ := ioutil.ReadFile(path)
	assert.NotContains(string(dat), "InsertNode")

	conn.InsertNode(&runtime.Node{
		Statistics: &data.Statistics{},
	})

	dat, _ = ioutil.ReadFile(path)
	assert.Contains(string(dat), "InsertNode")

	assert.NotContains(string(dat), "InsertLink")
	conn.InsertLink(&runtime.Link{}, time.Now())
	dat, _ = ioutil.ReadFile(path)
	assert.Contains(string(dat), "InsertLink")

	assert.NotContains(string(dat), "InsertGlobals")
	conn.InsertGlobals(&runtime.GlobalStats{}, time.Now(), runtime.GLOBAL_SITE, runtime.GLOBAL_DOMAIN)
	dat, _ = ioutil.ReadFile(path)
	assert.Contains(string(dat), "InsertGlobals")

	assert.NotContains(string(dat), "PruneNodes")
	conn.PruneNodes(time.Second)
	dat, _ = ioutil.ReadFile(path)
	assert.Contains(string(dat), "PruneNodes")

	assert.NotContains(string(dat), "Close")
	conn.Close()
	dat, _ = ioutil.ReadFile(path)
	assert.Contains(string(dat), "Close")

	os.Remove(path)
}
