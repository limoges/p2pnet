package msg

import (
	"encoding/binary"
	"fmt"
)

type AuthSessionHS1 struct {
	SessionID        uint32
	HandshakePayload []byte
}

func (m AuthSessionHS1) String() string {
	return fmt.Sprintf("AuthSessionHS1{SessionID:%v, Payload:%v}", m.SessionID,
		m.HandshakePayload)
}

func (m AuthSessionHS1) MinimumLength() int {
	return 4
}

func (m AuthSessionHS1) PayloadLength() int {
	return m.MinimumLength() + len(m.HandshakePayload)
}

func (m *AuthSessionHS1) UnmarshalBinary(data []byte) error {

	if len(data) < m.MinimumLength() {
		return ErrDataTooShort
	}

	m.SessionID = binary.BigEndian.Uint32(data[:4])
	m.HandshakePayload = make([]byte, len(data[4:]))
	copy(m.HandshakePayload, data[4:])
	return nil
}

func (m AuthSessionHS1) MarshalBinary() (data []byte, err error) {

	payloadBuf := make([]byte, m.PayloadLength())

	binary.BigEndian.PutUint32(payloadBuf[:4], m.SessionID)
	copy(payloadBuf, m.HandshakePayload)
	return createMessage(AUTH_SESSION_HS1, payloadBuf), nil
}
