package database

import (
	"sync"
	"testing"
	"time"

	"github.com/FreifunkBremen/yanic/runtime"
	"github.com/stretchr/testify/assert"
)

type testConn struct {
	Connection
	countClose int
	countPrune int
	sync.Mutex
}

func (c *testConn) Close() {
	c.Lock()
	c.countClose++
	c.Unlock()
}
func (c *testConn) GetClose() int {
	c.Lock()
	defer c.Unlock()
	return c.countClose
}
func (c *testConn) PruneNodes(time.Duration) {
	c.Lock()
	c.countPrune++
	c.Unlock()
}
func (c *testConn) GetPruneNodes() int {
	c.Lock()
	defer c.Unlock()
	return c.countPrune
}

func TestStart(t *testing.T) {
	assert := assert.New(t)

	conn := &testConn{}
	config := &runtime.Config{
		Database: struct {
			DeleteInterval runtime.Duration `toml:"delete_interval"`
			DeleteAfter    runtime.Duration `toml:"delete_after"`
			Connection     map[string]interface{}
		}{
			DeleteInterval: runtime.Duration{Duration: time.Millisecond * 10},
		},
	}
	assert.Nil(quit)

	Start(conn, config)
	assert.NotNil(quit)

	assert.Equal(0, conn.GetPruneNodes())
	time.Sleep(time.Millisecond * 12)
	assert.Equal(1, conn.GetPruneNodes())

	assert.Equal(0, conn.GetClose())
	Close(conn)
	assert.NotNil(quit)
	assert.Equal(1, conn.GetClose())

	time.Sleep(time.Millisecond * 12) // to reach timer.Stop() line

}
