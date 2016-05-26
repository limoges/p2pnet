package msg

import (
	"encoding/binary"
	"fmt"
)

// Gossip Announce
type GossipAnnounce struct {
	TTL      uint8
	reserved uint8
	DataType uint16
	Data     []byte
}

func (m GossipAnnounce) String() string {
	return fmt.Sprintf("Gossip.Announce{TTL:%v, DataType:%v, Data:%v}", m.TTL, m.DataType, m.Data)
}

func (m GossipAnnounce) MinimumLength() int {
	return 4
}

func (m GossipAnnounce) PayloadLength() int {
	return m.MinimumLength() + len(m.Data)
}

func (m *GossipAnnounce) UnmarshalBinary(data []byte) error {

	if len(data) < m.MinimumLength() {
		return ErrDataTooShort
	}

	m.TTL = data[0]
	m.reserved = data[1]
	m.DataType = binary.BigEndian.Uint16(data[2:4])
	m.Data = data[4:]
	return nil
}

func (m GossipAnnounce) MarshalBinary() (data []byte, err error) {

	payloadBuf := make([]byte, m.PayloadLength())
	payloadBuf[0] = m.TTL
	payloadBuf[1] = m.reserved
	binary.BigEndian.PutUint16(payloadBuf[2:4], m.DataType)
	copy(payloadBuf[4:], m.Data)
	return createMessage(GOSSIP_ANNOUNCE, payloadBuf), nil
}

// Gossip Notify
type GossipNotify struct {
	reserved uint16
	DataType uint16
}

func (m GossipNotify) String() string {
	return fmt.Sprintf("Gossip.Notify{DataType:%v}", m.DataType)
}

func (m GossipNotify) MinimumLength() int {
	return 4
}

func (m GossipNotify) PayloadLength() int {
	return m.MinimumLength()
}

func (m *GossipNotify) UnmarshalBinary(data []byte) error {

	if len(data) < m.MinimumLength() {
		return ErrDataTooShort
	}

	m.reserved = binary.BigEndian.Uint16(data[:2])
	m.DataType = binary.BigEndian.Uint16(data[2:4])
	return nil
}

func (m GossipNotify) MarshalBinary() (data []byte, err error) {

	payloadBuf := make([]byte, m.PayloadLength())
	binary.BigEndian.PutUint16(payloadBuf[:2], m.reserved)
	binary.BigEndian.PutUint16(payloadBuf[2:4], m.DataType)
	return createMessage(GOSSIP_NOTIFY, payloadBuf), nil
}

// Gossip Notification
type GossipNotification struct {
	HeaderID uint16
	DataType uint16
	Data     []byte
}

func (m GossipNotification) String() string {
	return fmt.Sprintf("Gossip.Notification{HeaderID:%v, DataType:%v, Data:%v}", m.HeaderID, m.DataType, m.Data)
}

func (m GossipNotification) MinimumLength() int {
	return 4
}

func (m GossipNotification) PayloadLength() int {
	return m.MinimumLength() + len(m.Data)
}

func (m *GossipNotification) UnmarshalBinary(data []byte) error {

	if len(data) < m.MinimumLength() {
		return ErrDataTooShort
	}

	m.HeaderID = binary.BigEndian.Uint16(data[:2])
	m.DataType = binary.BigEndian.Uint16(data[2:4])
	m.Data = data[4:]
	return nil
}

func (m GossipNotification) MarshalBinary() (data []byte, err error) {

	payloadBuf := make([]byte, m.PayloadLength())
	binary.BigEndian.PutUint16(payloadBuf[:2], m.HeaderID)
	binary.BigEndian.PutUint16(payloadBuf[2:4], m.DataType)
	copy(payloadBuf[4:], m.Data)
	return createMessage(GOSSIP_NOTIFICATION, payloadBuf), nil
}

type GossipValidation struct {
	MessageID uint16
	reserved  uint16
	Valid     bool
}

func (m *GossipValidation) String() string {
	return fmt.Sprintf("Gossip.Validation{MessageID:%v, Valid:%v}", m.MessageID, m.Valid)
}

func (m *GossipValidation) MinimumLength() int {
	return 4
}

func (m *GossipValidation) PayloadLength() int {
	return m.MinimumLength()
}

func (m *GossipValidation) UnmarshalBinary(data []byte) error {

	if len(data) < m.MinimumLength() {
		return ErrDataTooShort
	}

	m.MessageID = binary.BigEndian.Uint16(data[:2])
	m.reserved = binary.BigEndian.Uint16(data[2:4])
	m.Valid = m.reserved == 0x1
	return nil
}

func (m *GossipValidation) MarshalBinary() (data []byte, err error) {

	payloadBuf := make([]byte, m.PayloadLength())
	binary.BigEndian.PutUint16(payloadBuf[0:2], m.MessageID)
	if m.Valid {
		m.reserved = 0x1
	} else {
		m.reserved = 0
	}
	binary.BigEndian.PutUint16(payloadBuf[2:4], m.reserved)
	return createMessage(GOSSIP_VALIDATION, payloadBuf), nil
}
