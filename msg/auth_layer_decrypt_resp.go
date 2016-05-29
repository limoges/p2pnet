package msg

import (
	"fmt"
)

type AuthLayerDecryptResp struct {
}

func (m AuthLayerDecryptResp) String() string {
	return fmt.Sprintf("")
}

func (m AuthLayerDecryptResp) MinimumLength() int {
	return 0
}

func (m AuthLayerDecryptResp) PayloadLength() int {

	// TODO: Implementation
	return m.MinimumLength()
}

func (m *AuthLayerDecryptResp) UnmarshalBinary(data []byte) error {

	if len(data) < m.MinimumLength() {
		return ErrDataTooShort
	}

	// TODO: Implementation

	return nil
}

func (m AuthLayerDecryptResp) MarshalBinary() (data []byte, err error) {

	payloadBuf := make([]byte, m.PayloadLength())

	// TODO: Implementation

	return createMessage(AUTH_LAYER_DECRYPT_RESP, payloadBuf), nil
}
