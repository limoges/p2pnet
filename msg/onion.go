package msg

const (
	ONION_TUNNEL_BUILD    = 560
	ONION_TUNNEL_READY    = 561
	ONION_TUNNEL_INCOMING = 562
	ONION_TUNNEL_DESTROY  = 563
	ONION_TUNNEL_DATA     = 564
	ONION_ERROR           = 565
	ONION_COVER           = 566
	// Reserved up to 599.
)

type OnionTunnelIncoming struct {
	TunnelID           uint32
	SourceHostKeyInDER []byte
}

type OnionTunnelDestroy struct {
	TunnelID uint32
}

type OnionTunnelData struct {
	TunnelID uint32
	Data     []byte
}

type OnionError struct {
	RequestType uint16
	reserved    uint16
	TunnelID    uint32
}

type OnionCover struct {
	CoverSize uint16
	reserved  uint16
}
