package msg

const (
	AUTH_HANDSHAKE1 = 700
	AUTH_HANDSHAKE2 = 701
)

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
