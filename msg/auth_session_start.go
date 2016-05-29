package msg

import (
	"fmt"
)

type AuthSessionStart struct {
	Hostkey []byte
}

func (m AuthSessionStart) String() string {
	return fmt.Sprintf("AuthSessionStart{HopHostkey:%v}", m.Hostkey)
}

func (m AuthSessionStart) MinimumLength() int {
	return 0
}

func (m AuthSessionStart) PayloadLength() int {
	return m.MinimumLength() + len(m.Hostkey)
}

func (m *AuthSessionStart) UnmarshalBinary(data []byte) error {

	if len(data) < m.MinimumLength() {
		return ErrDataTooShort
	}

	m.Hostkey = make([]byte, len(data))
	copy(m.Hostkey, data)
	return nil
}

func (m AuthSessionStart) MarshalBinary() (data []byte, err error) {

	payloadBuf := make([]byte, m.PayloadLength())
	copy(payloadBuf, m.Hostkey)
	return createMessage(AUTH_SESSION_START, payloadBuf), nil
}
