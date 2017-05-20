package output

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	assert := assert.New(t)
	assert.Len(Adapters, 0)

	RegisterAdapter("blub", func(config map[string]interface{}) (Output, error) {
		return nil, nil
	})

	assert.Len(Adapters, 1)
}
