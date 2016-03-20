package jsontime

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMarshalTime(t *testing.T) {
	assert := assert.New(t)

	nativeTime, err := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
	assert.Nil(err)

	json, err := Time(nativeTime).MarshalJSON()
	assert.Nil(err)

	assert.Equal(`"2012-11-01T22:08:41"`, string(json))
}

func TestUnmarshalTime(t *testing.T) {
	assert := assert.New(t)
	jsonTime := Time{}

	err := json.Unmarshal([]byte(`"2012-11-01T22:08:41"`), &jsonTime)
	assert.Nil(err)
}
