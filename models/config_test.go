package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadConfig(t *testing.T) {
	assert := assert.New(t)

	config := ReadConfigFile("../config_example.yml")
	assert.NotNil(config)
}
