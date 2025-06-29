package geojson

import (
	"fmt"
	"os"
	"testing"

	"github.com/FreifunkBremen/yanic/runtime"
	"github.com/stretchr/testify/assert"
)

func TestOutput(t *testing.T) {
	assert := assert.New(t)

	out, err := Register(map[string]interface{}{})
	assert.Error(err)
	assert.Nil(out)

	out, err = Register(map[string]interface{}{
		"path": "/tmp/nodes.geojson",
	})
	if err := os.Remove("/tmp/nodes.geojson"); err != nil {
		fmt.Printf("during cleanup: %s\n", err)
	}
	assert.NoError(err)
	assert.NotNil(out)

	out.Save(&runtime.Nodes{})
	_, err = os.Stat("/tmp/nodes.geojson")
	assert.NoError(err)
}
