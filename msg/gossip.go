package msg

import (
	"bytes"
	"encoding/binary"
	"io"
)

const (
	GOSSIP_ANNOUNCE     = 500
	GOSSIP_NOTIFY       = 501
	GOSSIP_NOTIFICATION = 502
	GOSSIP_VALIDATION   = 503
	// Reserved up to 519.
)

type GossipAnnounce struct {
	TTL      uint8
	Reserved uint8
	DataType uint16
	Data     []byte
}

func (m GossipAnnounce) TypeId() uint16 {
	return GOSSIP_ANNOUNCE
}

func NewGossipAnnounce(data []byte) (GossipAnnounce, error) {

	var m GossipAnnounce
	var reader *bytes.Reader
	var err error

	m = GossipAnnounce{}
	reader = bytes.NewReader(data)

	if err = binary.Read(reader, binary.BigEndian, &m.TTL); err != nil {
		return m, err
	}
	if err = binary.Read(reader, binary.BigEndian, &m.Reserved); err != nil {
		return m, err
	}
	if err = binary.Read(reader, binary.BigEndian, &m.DataType); err != nil {
		return m, err
	}

	m.Data = make([]byte, reader.Len())
	if _, err = io.ReadFull(reader, m.Data); err != nil {
		return m, err
	}

	return m, nil
}

type GossipNotify struct {
	Reserved uint16
	DataType uint16
}

func (m GossipNotify) TypeId() uint16 {
	return GOSSIP_NOTIFY
}

func NewGossipNotify(data []byte) (GossipNotify, error) {

	var m GossipNotify
	var reader *bytes.Reader
	var err error

	m = GossipNotify{}
	reader = bytes.NewReader(data)

	if err = binary.Read(reader, binary.BigEndian, &m.Reserved); err != nil {
		return m, err
	}
	if err = binary.Read(reader, binary.BigEndian, &m.DataType); err != nil {
		return m, err
	}

	return m, nil
}

type GossipNotification struct {
	HeaderId uint16
	DataType uint16
	Data     []byte
}

func (m GossipNotification) TypeId() uint16 {
	return GOSSIP_NOTIFICATION
}

func NewGossipNotification(data []byte) (GossipNotification, error) {

	var m GossipNotification
	var reader *bytes.Reader
	var err error

	m = GossipNotification{}
	reader = bytes.NewReader(data)

	if err = binary.Read(reader, binary.BigEndian, &m.HeaderId); err != nil {
		return m, err
	}
	if err = binary.Read(reader, binary.BigEndian, &m.DataType); err != nil {
		return m, err
	}

	m.Data = make([]byte, reader.Len())
	if _, err = io.ReadFull(reader, m.Data); err != nil {
		return m, err
	}

	return m, nil
}

type GossipValidation struct {
	MessageId uint16
	Reserved  uint16
}

func (m GossipValidation) Valid() bool {
	return true
}

func (m *GossipValidation) SetValid(valid bool) {
	if valid {
		m.Reserved = 0x1
	} else {
		m.Reserved = 0x0
	}
}

func (m GossipValidation) TypeId() uint16 {
	return GOSSIP_VALIDATION
}

func NewGossipValidation(data []byte) (GossipValidation, error) {

	var m GossipValidation
	var reader *bytes.Reader
	var err error

	m = GossipValidation{}
	reader = bytes.NewReader(data)

	if err = binary.Read(reader, binary.BigEndian, &m.MessageId); err != nil {
		return m, err
	}
	if err = binary.Read(reader, binary.BigEndian, &m.Reserved); err != nil {
		return m, err
	}

	return m, nil
}
