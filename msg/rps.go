package msg

import (
	"bytes"
	"encoding/binary"
	"io"
)

const (
	RPS_QUERY = 540
	RPS_PEER  = 541
	// Reserved up to 559.
)

type RPSQuery struct {
	// This is empty.
}

func (m RPSQuery) TypeId() uint16 {
	return RPS_QUERY
}

func NewRPSQuery(data []byte) (RPSQuery, error) {
	return RPSQuery{}, nil
}

type RPSPeer struct {
	Port     uint16
	Reserved uint16
	IPAddr   [16]byte
	Hostkey  []byte
}

func (m RPSPeer) TypeId() uint16 {
	return RPS_PEER
}

func NewRPSPeer(data []byte) (RPSPeer, error) {

	var m RPSPeer
	var err error
	var reader *bytes.Reader

	m = RPSPeer{}
	reader = new(bytes.Reader)

	if err = binary.Read(reader, binary.BigEndian, &m.Port); err != nil {
		return m, err
	}
	if err = binary.Read(reader, binary.BigEndian, &m.Reserved); err != nil {
		return m, err
	}
	if _, err = io.ReadFull(reader, m.IPAddr[:]); err != nil {
		return m, err
	}
	m.Hostkey = make([]byte, reader.Len())
	if _, err = io.ReadFull(reader, m.Hostkey); err != nil {
		return m, err
	}
	return m, nil
}
