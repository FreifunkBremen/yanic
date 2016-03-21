package jsontime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNow(t *testing.T) {
	assert := assert.New(t)

	t1 := time.Now()
	t2 := Now()

	assert.InDelta(t1.Unix(), t2.Unix(), 1)
}

func TestMarshalTime(t *testing.T) {
	assert := assert.New(t)

	nativeTime, err := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
	assert.Nil(err)

	json, err := Time{nativeTime}.MarshalJSON()
	assert.Nil(err)

	assert.Equal(`"2012-11-01T22:08:41"`, string(json))
}

func TestUnmarshalTime(t *testing.T) {
	assert := assert.New(t)
	jsonTime := Time{}

	err := jsonTime.UnmarshalJSON([]byte(`"2012-11-01T22:08:41"`))
	assert.Nil(err)
	assert.False(jsonTime.IsZero())
}
