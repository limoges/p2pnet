package msg

import (
	"encoding/binary"
	"errors"
	"fmt"
)

const (
	HeaderLength = 4
)

const (
	IPLength = 16
)

const (
	GOSSIP_ANNOUNCE     = 500
	GOSSIP_NOTIFY       = 501
	GOSSIP_NOTIFICATION = 502
	GOSSIP_VALIDATION   = 503
	// Reserved up to 519.
	NSE_QUERY    = 520
	NSE_ESTIMATE = 521
	// Reserved up to 539.
	RPS_QUERY = 540
	RPS_PEER  = 541
	// Reserved up to 559.
	ONION_TUNNEL_BUILD    = 560
	ONION_TUNNEL_READY    = 561
	ONION_TUNNEL_INCOMING = 562
	ONION_TUNNEL_DESTROY  = 563
	ONION_TUNNEL_DATA     = 564
	ONION_ERROR           = 565
	ONION_COVER           = 566
	// Reserved up to 599.
	AUTH_SESSION_START        = 600
	AUTH_SESSION_HS1          = 601
	AUTH_SESSION_INCOMING_HS1 = 602
	AUTH_SESSION_HS2          = 603
	AUTH_SESSION_INCOMING_HS2 = 604
	AUTH_LAYER_ENCRYPT        = 605
	AUTH_LAYER_ENCRYPT_RESP   = 606
	AUTH_LAYER_DECRYPT        = 607
	AUTH_LAYER_DECRYPT_RESP   = 608
	AUTH_SESSION_CLOSE        = 609
	// Reserved up to 649.
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

// Header represents the protocol's common message header format.
// It is defined to contain message's size in the first 2 bytes,
// followed by the message's type in the following 2 bytes.
type Header struct {
	Size uint16
	Type uint16
}

func (h Header) String() string {
	return fmt.Sprintf("Header{Size:%v, Type:%v}", h.Size, h.Type)
}

func (h Header) MinimumLength() int {
	return HeaderLength
}

func (h Header) PayloadSize() int {
	return int(h.Size) - HeaderLength
}

func (h *Header) UnmarshalBinary(data []byte) error {

	if len(data) < h.MinimumLength() {
		return ErrDataTooShort
	}

	h.Size = binary.BigEndian.Uint16(data[:2])
	h.Type = binary.BigEndian.Uint16(data[2:4])
	return nil
}

func (h Header) MarshalBinary() (data []byte, err error) {

	data = make([]byte, HeaderLength)

	binary.BigEndian.PutUint16(data[0:2], h.Size)
	binary.BigEndian.PutUint16(data[2:4], h.Type)

	return data, nil
}

type NSEQuery struct {
	// This is empty.
}

func (m NSEQuery) String() string {
	return fmt.Sprintf("NSEQuery{}")
}

func (m *NSEQuery) MinimumLength() int {
	return 0
}

func (m *NSEQuery) PayloadLength() int {
	return m.MinimumLength()
}

func (m *NSEQuery) UnmarshalBinary(data []byte) error {
	return nil
}

func (m *NSEQuery) MarshalBinary() (data []byte, err error) {

	payloadBuf := make([]byte, 0)
	return createMessage(NSE_QUERY, payloadBuf), nil
}

type NSEEstimate struct {
	EstimatePeers        uint32
	EstimateStdDeviation uint32
}

func (m NSEEstimate) String() string {
	return fmt.Sprintf("NSEEstimate{EstimatePeers:%v, EstimateStdDeviation:%v}", m.EstimatePeers,
		m.EstimateStdDeviation)
}

func (m *NSEEstimate) MinimumLength() int {
	return 8
}

func (m *NSEEstimate) PayloadLength() int {
	return m.MinimumLength()
}

func (m *NSEEstimate) UnmarshalBinary(data []byte) error {

	if len(data) < m.MinimumLength() {
		return ErrDataTooShort
	}

	m.EstimatePeers = binary.BigEndian.Uint32(data[:4])
	m.EstimateStdDeviation = binary.BigEndian.Uint32(data[4:8])
	return nil
}

func (m *NSEEstimate) MarshalBinary() (data []byte, err error) {

	payloadBuf := make([]byte, m.PayloadLength())
	binary.BigEndian.PutUint32(payloadBuf[:4], m.EstimatePeers)
	binary.BigEndian.PutUint32(payloadBuf[4:8], m.EstimateStdDeviation)
	return createMessage(NSE_ESTIMATE, payloadBuf), nil
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
	// case ONION_TUNNEL_BUILD:
	// case ONION_TUNNEL_READY:
	// case ONION_TUNNEL_INCOMING:
	// case ONION_TUNNEL_DESTROY:
	// case ONION_TUNNEL_DATA:
	// case ONION_ERROR:
	// case ONION_COVER:
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
