package msg

import (
	"bytes"
	"encoding/binary"
	"io"
)

type RPSPeer struct {
	Port     uint16
	Reserved uint16
	IPAddr   []byte
	Hostkey  []byte
}

func (m RPSPeer) TypeId() uint16 {
	return RPS_PEER
}

func NewRPSPeer(data []byte) (RPSPeer, error) {

	m := RPSPeer{}
	reader := bytes.NewReader(data)
	if err := binary.Read(reader, binary.BigEndian, &m.Port); err != nil {
		return m, err
	}

	if err := binary.Read(reader, binary.BigEndian, &m.Reserved); err != nil {
		return m, err
	}

	m.IPAddr = make([]byte, 16)
	if _, err := io.ReadFull(reader, m.IPAddr); err != nil {
		return m, err
	}

	mustRead := len(data) - 16 - 4
	m.Hostkey = make([]byte, mustRead)
	if _, err := io.ReadFull(reader, m.Hostkey); err != nil {
		return m, err
	}

	return m, nil
}
