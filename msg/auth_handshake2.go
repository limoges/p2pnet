package msg

type AuthHandshake2 struct {
	Cipher []byte
}

func (m AuthHandshake2) TypeId() uint16 {
	return AUTH_HANDSHAKE2
}

func NewAuthHandshake2(data []byte) (AuthHandshake2, error) {

	var m AuthHandshake2

	m = AuthHandshake2{}
	m.Cipher = make([]byte, len(data))
	copy(m.Cipher, data)
	return m, nil
}
