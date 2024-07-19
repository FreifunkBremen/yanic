package runtime

type LinkType int

const (
	UnknownLinkType LinkType = iota
	WirelessLinkType
	TunnelLinkType
	OtherLinkType
)

func (lt LinkType) String() string {
	switch lt {
	case WirelessLinkType:
		return "wifi"
	case TunnelLinkType:
		return "vpn"
	case OtherLinkType:
		return "other"
	}
	return "unknown"
}

type LinkProtocol int

const (
	UnknownLinkProtocol LinkProtocol = iota
	BatadvLinkProtocol
	BabelLinkProtocol
	LLDPLinkProtocol
)

func (lp LinkProtocol) String() string {
	switch lp {
	case BatadvLinkProtocol:
		return "batadv"
	case BabelLinkProtocol:
		return "babel"
	case LLDPLinkProtocol:
		return "lldp"
	}
	return "unkown"
}

// Link represents a link between two nodes
type Link struct {
	SourceID       string
	SourceHostname string
	SourceAddress  string
	TargetID       string
	TargetAddress  string
	TargetHostname string
	TQ             float32
	Type           LinkType
	Protocol       LinkProtocol
}
