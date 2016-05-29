package msg

import (
	"encoding/binary"
	"fmt"
)

type AuthSessionIncomingHS1 struct {
	Reserved         uint16
	HostkeySize      uint16
	SourceHostkey    []byte
	HandshakePayload []byte
}

func (m AuthSessionIncomingHS1) String() string {
	return fmt.Sprintf(
		"AuthSessionIncomingHS1{HostkeySize:%v, SourceHostkey:%v, HandshakePayload%v",
		m.HostkeySize,
		m.SourceHostkey,
		m.HandshakePayload)
}

func (m AuthSessionIncomingHS1) MinimumLength() int {
	return 4
}

func (m AuthSessionIncomingHS1) PayloadLength() int {
	return m.MinimumLength() + len(m.SourceHostkey) + len(m.HandshakePayload)
}

func (m *AuthSessionIncomingHS1) UnmarshalBinary(data []byte) error {

	if len(data) < m.MinimumLength() {
		return ErrDataTooShort
	}

	m.HostkeySize = binary.BigEndian.Uint16(data[2:4])
	m.SourceHostkey = make([]byte, m.HostkeySize)
	copy(m.SourceHostkey, data[4:4+m.HostkeySize])
	m.HandshakePayload = make([]byte, len(data[4+m.HostkeySize:]))
	copy(m.HandshakePayload, data[4+m.HostkeySize:])
	return nil
}

func (m AuthSessionIncomingHS1) MarshalBinary() (data []byte, err error) {

	payloadBuf := make([]byte, m.PayloadLength())

	binary.BigEndian.PutUint16(data[2:4], m.HostkeySize)
	copy(payloadBuf[4:], m.SourceHostkey)
	copy(payloadBuf[4+len(m.SourceHostkey):], m.HandshakePayload)

	return createMessage(AUTH_SESSION_INCOMING_HS1, payloadBuf), nil
}
