package msg

import (
	"bytes"
	"encoding/binary"
)

type AuthSessionIncomingHS2 struct {
	SessionId uint32
	Payload   []byte
}

func (m AuthSessionIncomingHS2) TypeId() uint16 {
	return AUTH_SESSION_INCOMING_HS2
}

func NewAuthSessionIncomingHS2(data []byte) (AuthSessionIncomingHS2, error) {

	var m AuthSessionIncomingHS2
	var reader *bytes.Reader
	var err error

	m = AuthSessionIncomingHS2{}
	reader = bytes.NewReader(data)

	if err = binary.Read(reader, binary.BigEndian, &m.SessionId); err != nil {
		return m, err
	}
	data = data[4:]

	m.Payload = make([]byte, len(data))
	copy(m.Payload, data)

	return m, nil
}
