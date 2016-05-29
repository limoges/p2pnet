package msg

import (
	"fmt"
)

type AuthSessionIncomingHS2 struct {
}

func (m AuthSessionIncomingHS2) String() string {
	return fmt.Sprintf("")
}

func (m AuthSessionIncomingHS2) MinimumLength() int {
	return 0
}

func (m AuthSessionIncomingHS2) PayloadLength() int {

	// TODO: Implementation
	return m.MinimumLength()
}

func (m *AuthSessionIncomingHS2) UnmarshalBinary(data []byte) error {

	if len(data) < m.MinimumLength() {
		return ErrDataTooShort
	}

	// TODO: Implementation

	return nil
}

func (m AuthSessionIncomingHS2) MarshalBinary() (data []byte, err error) {

	payloadBuf := make([]byte, m.PayloadLength())

	// TODO: Implementation

	return createMessage(AUTH_SESSION_INCOMING_HS2, payloadBuf), nil
}
