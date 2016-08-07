package msg

import (
	"bytes"
	"encoding/binary"
	"io"
)

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
