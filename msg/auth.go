package msg

import (
	"encoding/binary"
	"fmt"
)

const (
	AUTH_SESSION_START        = 600
	AUTH_SESSION_HS1          = 601
	AUTH_SESSION_INCOMING_HS1 = 602
	AUTH_SESSION_HS2          = 603
	AUTH_SESSION_INCOMING_HS2 = 604
	AUTH_LAYER_ENCRYPT        = 605
	AUTH_LAYER_ENCRYPT_RESP   = 606
	AUTH_LAYER_DECRYPT        = 607
	AUTH_LAYER_DECRYPT_RESP   = 608
	AUTH_SESSION_CLOSE        = 609
	// Reserved up to 649.
)

type AuthSessionStart struct {
	HopHostKey []byte
}

func (m AuthSessionStart) String() string {
	return fmt.Sprintf("AuthSessionStart{HopHostKey:%v}", m.HopHostKey)
}

func (m AuthSessionStart) MinimumLength() int {
	return 0
}

func (m AuthSessionStart) PayloadLength() int {
	return m.MinimumLength() + len(m.HopHostKey)
}

func (m *AuthSessionStart) UnmarshalBinary(data []byte) error {

	if len(data) < m.MinimumLength() {
		return ErrDataTooShort
	}

	m.HopHostKey = make([]byte, len(data))
	copy(m.HopHostKey, data)
	return nil
}

func (m AuthSessionStart) MarshalBinary() (data []byte, err error) {

	payloadBuf := make([]byte, m.PayloadLength())
	copy(payloadBuf, m.HopHostKey)
	return createMessage(AUTH_SESSION_START, payloadBuf), nil
}

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
	return createMessage(AUTH_SESSION_START, payloadBuf), nil
}
