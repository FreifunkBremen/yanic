package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNow(t *testing.T) {
	assert := assert.New(t)
	assert.Len(Adapters, 0)

	RegisterAdapter("blub", func(config interface{}) (Connection, error) {
		return nil, nil
	})

	assert.Len(Adapters, 1)
}
