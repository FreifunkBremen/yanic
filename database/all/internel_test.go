package all

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/FreifunkBremen/yanic/database"
	"github.com/FreifunkBremen/yanic/runtime"
	"github.com/stretchr/testify/assert"
)

type testConn struct {
	database.Connection
	countNode    int
	countLink    int
	countGlobals int
	countPrune   int
	countClose   int
	sync.Mutex
}

func (c *testConn) InsertNode(node *runtime.Node) {
	c.Lock()
	c.countNode++
	c.Unlock()
}
func (c *testConn) GetNode() int {
	c.Lock()
	defer c.Unlock()
	return c.countNode
}
func (c *testConn) InsertLink(link *runtime.Link, time time.Time) {
	c.Lock()
	c.countLink++
	c.Unlock()
}
func (c *testConn) GetLink() int {
	c.Lock()
	defer c.Unlock()
	return c.countLink
}
func (c *testConn) InsertGlobals(stats *runtime.GlobalStats, time time.Time, site string) {
	c.Lock()
	c.countGlobals++
	c.Unlock()
}
func (c *testConn) GetGlobal() int {
	c.Lock()
	defer c.Unlock()
	return c.countGlobals
}
func (c *testConn) PruneNodes(time.Duration) {
	c.Lock()
	c.countPrune++
	c.Unlock()
}
func (c *testConn) GetPrune() int {
	c.Lock()
	defer c.Unlock()
	return c.countPrune
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

func TestStart(t *testing.T) {
	assert := assert.New(t)

	globalConn := &testConn{}
	database.RegisterAdapter("a", func(config map[string]interface{}) (database.Connection, error) {
		return globalConn, nil
	})
	database.RegisterAdapter("b", func(config map[string]interface{}) (database.Connection, error) {
		return globalConn, nil
	})
	database.RegisterAdapter("c", func(config map[string]interface{}) (database.Connection, error) {
		return globalConn, nil
	})
	database.RegisterAdapter("d", func(config map[string]interface{}) (database.Connection, error) {
		return nil, nil
	})
	database.RegisterAdapter("e", func(config map[string]interface{}) (database.Connection, error) {
		return nil, errors.New("blub")
	})
	allConn, err := Connect(map[string]interface{}{
		"a": []map[string]interface{}{
			map[string]interface{}{
				"enable": false,
				"path":   "a1",
			},
			map[string]interface{}{
				"path": "a2",
			},
			map[string]interface{}{
				"enable": true,
				"path":   "a3",
			},
		},
		"b": nil,
		"c": []map[string]interface{}{
			map[string]interface{}{
				"path": "c1",
			},
		},
		// fetch continue command in Connect
		"d": []map[string]interface{}{
			map[string]interface{}{
				"path": "d0",
			},
		},
	})
	assert.NoError(err)

	assert.Equal(0, globalConn.GetNode())
	allConn.InsertNode(nil)
	assert.Equal(3, globalConn.GetNode())

	assert.Equal(0, globalConn.GetLink())
	allConn.InsertLink(nil, time.Now())
	assert.Equal(3, globalConn.GetLink())

	assert.Equal(0, globalConn.GetGlobal())
	allConn.InsertGlobals(nil, time.Now(), runtime.GLOBAL_SITE)
	assert.Equal(3, globalConn.GetGlobal())

	assert.Equal(0, globalConn.GetPrune())
	allConn.PruneNodes(time.Second)
	assert.Equal(3, globalConn.GetPrune())

	assert.Equal(0, globalConn.GetClose())
	allConn.Close()
	assert.Equal(3, globalConn.GetClose())

	_, err = Connect(map[string]interface{}{
		"e": []map[string]interface{}{
			map[string]interface{}{},
		},
	})
	assert.Error(err)

	// wrong format -> the only panic in Register
	assert.Panics(func() {
		Connect(map[string]interface{}{
			"e": true,
		})
	})
}
