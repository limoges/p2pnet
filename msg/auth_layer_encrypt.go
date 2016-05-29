package msg

import (
	"fmt"
)

type AuthLayerEncrypt struct {
}

func (m AuthLayerEncrypt) String() string {
	return fmt.Sprintf("")
}

func (m AuthLayerEncrypt) MinimumLength() int {
	return 0
}

func (m AuthLayerEncrypt) PayloadLength() int {

	// TODO: Implementation
	return m.MinimumLength()
}

func (m *AuthLayerEncrypt) UnmarshalBinary(data []byte) error {

	if len(data) < m.MinimumLength() {
		return ErrDataTooShort
	}

	// TODO: Implementation

	return nil
}

func (m AuthLayerEncrypt) MarshalBinary() (data []byte, err error) {

	payloadBuf := make([]byte, m.PayloadLength())

	// TODO: Implementation

	return createMessage(AUTH_LAYER_ENCRYPT, payloadBuf), nil
}
