package msg

import (
	"bytes"
	"encoding/binary"
	"io"
)

type AuthSessionIncomingHS1 struct {
	Reserved         uint16
	HostkeySize      uint16
	SourceHostkey    []byte
	HandshakePayload []byte
}

func (m AuthSessionIncomingHS1) TypeId() uint16 {
	return AUTH_SESSION_INCOMING_HS1
}

func NewAuthSessionIncomingHS1(data []byte) (AuthSessionIncomingHS1, error) {

	m := AuthSessionIncomingHS1{}
	reader := bytes.NewReader(data)
	if err := binary.Read(reader, binary.BigEndian, &m.Reserved); err != nil {
		return m, err
	}
	if err := binary.Read(reader, binary.BigEndian, &m.HostkeySize); err != nil {
		return m, err
	}

	mustRead := int(m.HostkeySize)
	hostkey := make([]byte, mustRead)
	if _, err := io.ReadFull(reader, hostkey); err != nil {
		return m, err
	}
	m.SourceHostkey = hostkey

	mustRead = len(data) - int(m.HostkeySize) - 4
	payload := make([]byte, mustRead)
	if _, err := io.ReadFull(reader, payload); err != nil {
		return m, err
	}
	m.HandshakePayload = payload

	return m, nil
}
