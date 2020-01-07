package respond

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRespone(t *testing.T) {
	assert := assert.New(t)

	addr, err := net.ResolveUDPAddr("udp6", "[fe80::2]:8080")
	assert.NoError(err)

	data, err := NewRespone(nil, addr)
	assert.NoError(err)
	assert.Equal("[fe80::2]:8080", data.Address.String())
	assert.Equal([]uint8{0xca, 0x2b, 0xcd, 0xc9, 0xe1, 0x2, 0x0, 0x0, 0x0, 0xff, 0xff}, data.Raw)
}
