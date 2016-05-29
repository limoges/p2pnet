package msg

import (
	"bytes"
	"encoding/binary"
	"io"
)

type OnionTunnelReady struct {
	TunnelId   uint32
	DstHostkey []byte
}

func (m OnionTunnelReady) TypeId() uint16 {
	return ONION_TUNNEL_READY
}

func NewOnionTunnelReady(data []byte) (OnionTunnelReady, error) {

	m := OnionTunnelReady{}
	reader := bytes.NewReader(data)

	// TunnelID field
	if err := binary.Read(reader, binary.BigEndian, &m.TunnelId); err != nil {
		return m, err
	}

	// Destination Hostkey
	mustRead := len(data) - 4
	m.DstHostkey = make([]byte, mustRead)
	if _, err := io.ReadFull(reader, m.DstHostkey); err != nil {
		return m, err
	}

	return m, nil
}
