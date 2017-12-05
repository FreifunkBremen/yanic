package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNodeinfoBatAddresses(t *testing.T) {
	assert := assert.New(t)
	iface := &NetworkInterface{
		Interfaces: struct {
			Wireless []string `json:"wireless,omitempty"`
			Other    []string `json:"other,omitempty"`
			Tunnel   []string `json:"tunnel,omitempty"`
		}{
			Wireless: nil,
			Other:    []string{"aa:aa:aa:aa:aa", "aa:aa:aa:aa:ab"},
			Tunnel:   []string{},
		},
	}

	addr := iface.Addresses()
	assert.NotNil(addr)
	assert.Equal([]string{"aa:aa:aa:aa:aa", "aa:aa:aa:aa:ab"}, addr)
}
