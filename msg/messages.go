package msg

import (
	"errors"
	"fmt"
)

const (
	HeaderLength = 4
	IPLength     = 16
)

var (
	ErrDataTooShort = errors.New("unmarshal: data is shorter than expected")
)

func createMessage(messageType int, payloadBuf []byte) (data []byte) {

	header := Header{
		Size: uint16(len(payloadBuf) + HeaderLength),
		Type: uint16(messageType),
	}
	headerBuf, _ := header.MarshalBinary()
	data = append(data, headerBuf...)
	data = append(data, payloadBuf...)
	return data
}

type Message interface {
	MinimumLength() int
	PayloadLength() int
	MarshalBinary() (data []byte, err error)
	UnmarshalBinary(data []byte) error
}

func NewMessage(h Header, payload []byte) Message {

	var m Message
	switch h.Type {
	case GOSSIP_ANNOUNCE:
		m = &GossipAnnounce{}
	case GOSSIP_NOTIFY:
		m = &GossipNotify{}
	case GOSSIP_NOTIFICATION:
		m = &GossipNotification{}
	case GOSSIP_VALIDATION:
		m = &GossipValidation{}
	case NSE_QUERY:
		m = &NSEQuery{}
	case NSE_ESTIMATE:
		m = &NSEEstimate{}
	case RPS_QUERY:
		m = &RPSQuery{}
	case RPS_PEER:
		m = &RPSPeer{}
	case ONION_TUNNEL_BUILD:
		m = &OnionTunnelBuild{}
	case ONION_TUNNEL_READY:
		m = &OnionTunnelReady{}
	case ONION_TUNNEL_INCOMING:
		m = &OnionTunnelIncoming{}
	case ONION_TUNNEL_DESTROY:
		m = &OnionTunnelDestroy{}
	case ONION_TUNNEL_DATA:
		m = &OnionTunnelData{}
	case ONION_ERROR:
		m = &OnionError{}
	case ONION_COVER:
		m = &OnionCover{}
	// case AUTH_SESSION_START:
	// case AUTH_SESSION_HS1:
	// case AUTH_SESSION_INCOMING_HS1:
	// case AUTH_SESSION_HS2:
	// case AUTH_SESSION_INCOMING_HS2:
	// case AUTH_LAYER_ENCRYPT:
	// case AUTH_LAYER_ENCRYPT_RESP:
	// case AUTH_LAYER_DECRYPT:
	// case AUTH_LAYER_DECRYPT_RESP:
	// case AUTH_SESSION_CLOSE:
	default:
		fmt.Printf("Unhandled message type: %v\n", h.Type)
	}
	m.UnmarshalBinary(payload)
	fmt.Println(m)
	return m
}
