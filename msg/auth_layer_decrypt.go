package msg

import (
	_ "fmt"
)

// type AuthLayerDecrypt struct {
// 	NbLayers         uint8
// 	Reserved         uint8
// 	RequestId        uint16
// 	SessionIds       []uint32
// 	EncryptedPayload []byte
// }
//
// func (m AuthLayerDecrypt) String() string {
// 	return fmt.Sprintf("%v", m)
// }
//
// func (m AuthLayerDecrypt) HeaderLengthInBytes() int {
// 	return 4
// }
//
// func (m AuthLayerDecrypt) MinimumLength() int {
// 	return 4
// }
//
// func (m AuthLayerDecrypt) PayloadLength() int {
//
// 	return m.MinimumLength() + len(sessionIDs)*4 + len(m.EncryptedPayload)
// }
//
// func NewAuthLayerDecrypt(data []byte) (AuthLayerDecrypt, error) {
//
// 	m := AuthLayerDecrypt{}
//
// 	if len(data) < m.MinimumLength() {
// 		return ErrDataTooShort
// 	}
//
// 	// Read the header
// 	m.NbLayers = uint8(data[0])
// 	m.Reserved = uint8(data[1])
// 	m.RequestId = binary.BigEndian.Uint16(data[2:4])
//
// 	// Read the session ids
// 	data = data[4:]
// 	m.SessionIds = make([]uint32, m.NbLayers)
// 	for i := range m.NbLayers {
// 		m.SessionIds[i] = binary.BigEndian.Uint32(data[:4])
// 		data = data[4:]
// 	}
//
// 	// Read the payload
// 	m.EncryptedPayload = make([]byte, len(data))
// 	copy(m.EncryptedPayload, data)
//
// 	return m, nil
// }
//
// func (m AuthLayerDecrypt) MarshalBinary() (data []byte, err error) {
//
// 	buf := new(bytes.Buffer)
//
// 	buf[0] = m.NbLayers
// 	buf[1] = m.Reserved
// 	binary.BigEndian.PutUint16(buf[2:4], m.RequestId)
//
// 	buf = append(buf, binary.BigEndian.PutUint16(
//
// 	return createMessage(AUTH_LAYER_DECRYPT, buf), nil
// }
