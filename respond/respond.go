package respond

import (
	"net"
)

const (
	// default multicast group used by announced
	multicastAddressDefault = "ff02:0:0:0:0:0:2:1001"

	// default udp port used by announced
	port = 1001

	// maximum receivable size
	maxDataGramSize = 8192
)

// Response of the respond request
type Response struct {
	Address *net.UDPAddr
	Raw     []byte
}
