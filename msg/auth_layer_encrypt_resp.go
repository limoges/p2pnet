package msg

import (
	"fmt"
)

type AuthLayerEncryptResp struct {
}

func (m AuthLayerEncryptResp) String() string {
	return fmt.Sprintf("")
}

func (m AuthLayerEncryptResp) MinimumLength() int {
	return 0
}

func (m AuthLayerEncryptResp) PayloadLength() int {

	// TODO: Implementation
	return m.MinimumLength()
}

func (m *AuthLayerEncryptResp) UnmarshalBinary(data []byte) error {

	if len(data) < m.MinimumLength() {
		return ErrDataTooShort
	}

	// TODO: Implementation

	return nil
}

func (m AuthLayerEncryptResp) MarshalBinary() (data []byte, err error) {

	payloadBuf := make([]byte, m.PayloadLength())

	// TODO: Implementation

	return createMessage(AUTH_LAYER_ENCRYPT_RESP, payloadBuf), nil
}
