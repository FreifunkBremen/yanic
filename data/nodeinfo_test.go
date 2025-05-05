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

func TestNodeinfo(t *testing.T) {
	assert := assert.New(t)
	obj := &Nodeinfo{}
	testfile("nodeinfo.json", obj)

	assert.Equal("stable", obj.Software.Autoupdater.Branch)

	assert.Equal("gluon-v2016.1.2", obj.Software.Firmware.Base)
	assert.Equal("2016.1.2+bremen1", obj.Software.Firmware.Release)

	assert.Equal("TP-Link TL-WDR4900 v1", obj.Hardware.Model)
}
