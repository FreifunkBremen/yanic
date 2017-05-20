package output

import (
	"testing"
	"time"

	"github.com/FreifunkBremen/yanic/runtime"
	"github.com/stretchr/testify/assert"
)

type testConn struct {
	Output
	CountSave int
}

func (c *testConn) Save(nodes *runtime.Nodes) {
	c.CountSave++
}

func TestStart(t *testing.T) {
	assert := assert.New(t)

	conn := &testConn{}
	config := &runtime.Config{
		Nodes: struct {
			Enable       bool             `toml:"enable"`
			StatePath    string           `toml:"state_path"`
			SaveInterval runtime.Duration `toml:"save_interval"`
			OfflineAfter runtime.Duration `toml:"offline_after"`
			PruneAfter   runtime.Duration `toml:"prune_after"`
			Output       map[string]interface{}
		}{
			SaveInterval: runtime.Duration{Duration: time.Millisecond * 10},
		},
	}
	assert.Nil(quit)

	Start(conn, nil, config)
	assert.NotNil(quit)

	assert.Equal(0, conn.CountSave)
	time.Sleep(time.Millisecond * 12)
	assert.Equal(1, conn.CountSave)

	time.Sleep(time.Millisecond * 12)
	Close()
	assert.Equal(2, conn.CountSave)

}
