package meshviewer

import (
	"os"
	"testing"

	"yanic/runtime"
	"github.com/stretchr/testify/assert"
)

func TestOutput(t *testing.T) {
	assert := assert.New(t)

	// no version defined
	out, err := Register(map[string]interface{}{})
	assert.Error(err)
	assert.Nil(out)

	// no nodes path defined
	out, err = Register(map[string]interface{}{
		"version": int64(1),
	})
	assert.NoError(err)
	assert.NotNil(out)
	assert.Panics(func() {
		out.Save(&runtime.Nodes{})
	})

	out, err = Register(map[string]interface{}{
		"version":    int64(2),
		"nodes_path": "/tmp/nodes.json",
		"graph_path": "/tmp/graph.json",
	})
	os.Remove("/tmp/nodes.json")
	os.Remove("/tmp/graph.json")
	assert.NoError(err)
	assert.NotNil(out)

	out.Save(&runtime.Nodes{})
	_, err = os.Stat("/tmp/nodes.json")
	assert.NoError(err)
	_, err = os.Stat("/tmp/graph.json")
	assert.NoError(err)
}
