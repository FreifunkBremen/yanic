package all

import (
	"errors"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"yanic/output"
	"yanic/runtime"
)

type testOutput struct {
	output.Output
	countSave int
	sync.Mutex
}

func (c *testOutput) Save(nodes *runtime.Nodes) {
	c.Lock()
	c.countSave++
	c.Unlock()
}
func (c *testOutput) Get() int {
	c.Lock()
	defer c.Unlock()
	return c.countSave
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
			{
				"enable": false,
				"path":   "a1",
			},
			{
				"path": "a2",
			},
			{
				"enable": true,
				"path":   "a3",
			},
		},
		"b": nil,
		"c": []map[string]interface{}{
			{
				"path":   "c1",
				"filter": map[string]interface{}{},
			},
		},
		// fetch continue command in Connect
		"d": []map[string]interface{}{
			{
				"path": "d0",
			},
		},
	})
	assert.NoError(err)

	assert.Equal(0, globalOutput.Get())
	allOutput.Save(nodes)
	assert.Equal(3, globalOutput.Get())

	// wrong format - map
	_, err = Register(map[string]interface{}{
		"e": []interface{}{
			false,
		},
	})
	assert.Error(err)

	// wrong format - array
	_, err = Register(map[string]interface{}{
		"e": true,
	})
	assert.Error(err)

	// output error
	_, err = Register(map[string]interface{}{
		"e": []interface{}{
			map[string]interface{}{
				"enable": true,
			},
		},
	})
	assert.Error(err)

	// output error invalid config of filter
	_, err = Register(map[string]interface{}{
		"a": []interface{}{
			map[string]interface{}{
				"enable": true,
				"filter": map[string]interface{}{
					"blocklist": true,
				},
			},
		},
	})
	assert.Error(err)
}
