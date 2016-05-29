package msg

import (
	"fmt"
)

type AuthSessionClose struct {
}

func (m AuthSessionClose) String() string {
	return fmt.Sprintf("")
}

func (m AuthSessionClose) MinimumLength() int {
	return 0
}

func (m AuthSessionClose) PayloadLength() int {

	// TODO: Implementation
	return m.MinimumLength()
}

func (m *AuthSessionClose) UnmarshalBinary(data []byte) error {

	if len(data) < m.MinimumLength() {
		return ErrDataTooShort
	}

	// TODO: Implementation

	return nil
}

func (m AuthSessionClose) MarshalBinary() (data []byte, err error) {

	payloadBuf := make([]byte, m.PayloadLength())

	// TODO: Implementation

	return createMessage(AUTH_SESSION_CLOSE, payloadBuf), nil
}
