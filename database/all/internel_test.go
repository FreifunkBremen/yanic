package all

import (
	"errors"
	"testing"
	"time"

	"github.com/FreifunkBremen/yanic/database"
	"github.com/FreifunkBremen/yanic/runtime"
	"github.com/stretchr/testify/assert"
)

func TestStart(t *testing.T) {
	assert := assert.New(t)

	database.RegisterAdapter("d", func(config map[string]interface{}) (database.Connection, error) {
		return nil, nil
	})
	database.RegisterAdapter("e", func(config map[string]interface{}) (database.Connection, error) {
		return nil, errors.New("blub")
	})
	// Test for PruneNodes (by start)
	assert.Nil(quit)
	err := Start(runtime.DatabaseConfig{
		DeleteInterval: runtime.Duration{Duration: time.Millisecond},
		Connection: map[string]interface{}{
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
		},
	})
	assert.NoError(err)
	assert.NotNil(quit)

	_, err = Connect(map[string]interface{}{
		"e": []map[string]interface{}{
			map[string]interface{}{},
		},
	})
	assert.Error(err)

	// wrong format
	_, err = Connect(map[string]interface{}{
		"e": true,
	})
	assert.Error(err)
}
