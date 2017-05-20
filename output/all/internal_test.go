package all

import (
	"errors"
	"testing"

	"github.com/FreifunkBremen/yanic/output"
	"github.com/FreifunkBremen/yanic/runtime"
	"github.com/stretchr/testify/assert"
)

type testOutput struct {
	output.Output
	CountSave int
}

func (c *testOutput) Save(nodes *runtime.Nodes) {
	c.CountSave++
}

func TestStart(t *testing.T) {
	assert := assert.New(t)

	nodes := &runtime.Nodes{}

	globalOutput := &testOutput{}
	output.RegisterAdapter("a", func(config map[string]interface{}) (output.Output, error) {
		return globalOutput, nil
	})
	output.RegisterAdapter("b", func(config map[string]interface{}) (output.Output, error) {
		return globalOutput, nil
	})
	output.RegisterAdapter("c", func(config map[string]interface{}) (output.Output, error) {
		return globalOutput, nil
	})
	output.RegisterAdapter("d", func(config map[string]interface{}) (output.Output, error) {
		return nil, nil
	})
	output.RegisterAdapter("e", func(config map[string]interface{}) (output.Output, error) {
		return nil, errors.New("blub")
	})
	allOutput, err := Register(map[string]interface{}{
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
				"path":   "c1",
				"filter": map[string]interface{}{},
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

	assert.Equal(0, globalOutput.CountSave)
	allOutput.Save(nodes)
	assert.Equal(3, globalOutput.CountSave)

	_, err = Register(map[string]interface{}{
		"e": []map[string]interface{}{
			map[string]interface{}{},
		},
	})
	assert.Error(err)

	// wrong format -> the only panic in Register
	assert.Panics(func() {
		Register(map[string]interface{}{
			"e": true,
		})
	})
}
