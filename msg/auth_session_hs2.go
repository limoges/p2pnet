package msg

import (
	"fmt"
)

type AuthSessionHS2 struct {
}

func (m AuthSessionHS2) String() string {
	return fmt.Sprintf("")
}

func (m AuthSessionHS2) MinimumLength() int {
	return 0
}

func (m AuthSessionHS2) PayloadLength() int {

	// TODO: Implementation
	return m.MinimumLength()
}

func (m *AuthSessionHS2) UnmarshalBinary(data []byte) error {

	if len(data) < m.MinimumLength() {
		return ErrDataTooShort
	}

	// TODO: Implementation

	return nil
}

func (m AuthSessionHS2) MarshalBinary() (data []byte, err error) {

	payloadBuf := make([]byte, m.PayloadLength())

	// TODO: Implementation

	return createMessage(AUTH_SESSION_HS2, payloadBuf), nil
}
