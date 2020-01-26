package prometheus_sd

import (
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
		"path": "/tmp/prometheus_sd.json",
	})
	os.Remove("/tmp/prometheus_sd.json")
	assert.NoError(err)
	assert.NotNil(out)

	out.Save(&runtime.Nodes{})
	_, err = os.Stat("/tmp/prometheus_sd.json")
	assert.NoError(err)
}
