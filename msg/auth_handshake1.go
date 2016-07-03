package msg

type AuthHandshake1 struct {
	Cipher []byte
}

func (m AuthHandshake1) TypeId() uint16 {
	return AUTH_HANDSHAKE1
}

func NewAuthHandshake1(data []byte) (AuthHandshake1, error) {

	var m AuthHandshake1

	m = AuthHandshake1{}
	m.Cipher = make([]byte, len(data))
	copy(m.Cipher, data)
	return m, nil
}
