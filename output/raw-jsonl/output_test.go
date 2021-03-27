package jsonlines

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOutput(t *testing.T) {
	assert := assert.New(t)

	out, err := Register(map[string]interface{}{})
	assert.Error(err)
	assert.Nil(out)

	out, err = Register(map[string]interface{}{
		"path": "/tmp/raw.jsonl",
	})
	os.Remove("/tmp/raw.jsonl")
	assert.NoError(err)
	assert.NotNil(out)

	out.Save(createTestNodes())
	_, err = os.Stat("/tmp/raw.jsonl")
	assert.NoError(err)
}
