package meshviewerFFRGB

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"yanic/runtime"
)

func TestOutput(t *testing.T) {
	assert := assert.New(t)
	out, err := Register(map[string]interface{}{})
	assert.Error(err)
	assert.Nil(out)

	out, err = Register(map[string]interface{}{
		"path": "/tmp/meshviewer.json",
	})
	os.Remove("/tmp/meshviewer.json")
	assert.NoError(err)
	assert.NotNil(out)

	out.Save(&runtime.Nodes{})
	_, err = os.Stat("/tmp/meshviewer.json")
	assert.NoError(err)
}
