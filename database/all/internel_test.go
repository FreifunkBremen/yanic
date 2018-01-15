package all

import (
	"errors"
	"testing"
	"time"

	"github.com/FreifunkBremen/yanic/database"
	"github.com/FreifunkBremen/yanic/lib/duration"
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
	err := Start(database.Config{
		DeleteInterval: duration.Duration{Duration: time.Millisecond},
		Connection: map[string]interface{}{
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
			"c": []interface{}{
				map[string]interface{}{
					"path": "c1",
				},
			},
			// fetch continue command in Connect
			"d": []interface{}{
				map[string]interface{}{
					"path": "d0",
				},
			},
		},
	})
	assert.NoError(err)
	assert.NotNil(quit)

	// connection type not found
	_, err = Connect(map[string]interface{}{
		"e": []map[string]interface{}{
			{},
		},
	})
	assert.Error(err)

	// test close
	Close()

	// wrong format
	err = Start(database.Config{
		Connection: map[string]interface{}{
			"e": true,
		},
	})
	assert.Error(err)
}
