package msg

import (
	"bytes"
	"encoding/binary"
)

type AuthSessionHS1 struct {
	SessionId        uint32
	HandshakePayload []byte
}

func (m AuthSessionHS1) TypeId() uint16 {
	return AUTH_SESSION_HS1
}

func NewAuthSessionHS1(data []byte) (AuthSessionHS1, error) {

	m := AuthSessionHS1{}
	buf := bytes.NewReader(data)
	if err := binary.Read(buf, binary.BigEndian, &m.SessionId); err != nil {
		return m, err
	}
	data = data[4:]
	m.HandshakePayload = make([]byte, len(data))
	copy(m.HandshakePayload, data)
	return m, nil
}
