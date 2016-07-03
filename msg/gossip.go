package msg

const (
	GOSSIP_ANNOUNCE     = 500
	GOSSIP_NOTIFY       = 501
	GOSSIP_NOTIFICATION = 502
	GOSSIP_VALIDATION   = 503
	// Reserved up to 519.
)

// Gossip Announce
type GossipAnnounce struct {
	TTL      uint8
	reserved uint8
	DataType uint16
	Data     []byte
}

// Gossip Notify
type GossipNotify struct {
	reserved uint16
	DataType uint16
}

// Gossip Notification
type GossipNotification struct {
	HeaderID uint16
	DataType uint16
	Data     []byte
}

type GossipValidation struct {
	MessageID uint16
	reserved  uint16
	Valid     bool
}
