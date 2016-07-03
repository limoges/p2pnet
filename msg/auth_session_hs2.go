package msg

import (
	"bytes"
	"encoding/binary"
)

type AuthSessionHS2 struct {
	SessionId        uint32
	HandshakePayload []byte
}

func (m AuthSessionHS2) TypeId() uint16 {
	return AUTH_SESSION_HS2
}

func NewAuthSessionHS2(data []byte) (AuthSessionHS2, error) {

	var m AuthSessionHS2
	var reader *bytes.Reader
	var err error

	m = AuthSessionHS2{}
	reader = bytes.NewReader(data)
	if err = binary.Read(reader, binary.BigEndian, &m.SessionId); err != nil {
		return m, err
	}
	data = data[4:]
	m.HandshakePayload = make([]byte, len(data))
	copy(m.HandshakePayload, data)
	return m, nil
}
