package msg

import (
	"bytes"
	"encoding/binary"
	"io"
)

type OnionTunnelBuild struct {
	Reserved   uint16
	Port       uint16
	IPAddr     []byte
	DstHostkey []byte
}

func (m OnionTunnelBuild) TypeId() uint16 {
	return ONION_TUNNEL_BUILD
}

func NewOnionTunnelBuild(data []byte) (OnionTunnelBuild, error) {

	m := OnionTunnelBuild{}
	reader := bytes.NewReader(data)

	// Reserved field
	if err := binary.Read(reader, binary.BigEndian, &m.Reserved); err != nil {
		return m, err
	}

	// Port field
	if err := binary.Read(reader, binary.BigEndian, &m.Port); err != nil {
		return m, err
	}

	// IP Field
	m.IPAddr = make([]byte, 16)
	if _, err := io.ReadFull(reader, m.IPAddr); err != nil {
		return m, err
	}

	mustRead := len(data) - 16 - 4
	m.DstHostkey = make([]byte, mustRead)
	if _, err := io.ReadFull(reader, m.DstHostkey); err != nil {
		return m, err
	}

	return m, nil
}
