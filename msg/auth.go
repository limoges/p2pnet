package msg

import (
	"bytes"
	"encoding/binary"
	"io"
)

const (
	AUTH_SESSION_START        = 600
	AUTH_SESSION_HS1          = 601
	AUTH_SESSION_INCOMING_HS1 = 602
	AUTH_SESSION_HS2          = 603
	AUTH_SESSION_INCOMING_HS2 = 604
	AUTH_LAYER_ENCRYPT        = 605
	AUTH_LAYER_ENCRYPT_RESP   = 606
	AUTH_LAYER_DECRYPT        = 607
	AUTH_LAYER_DECRYPT_RESP   = 608
	AUTH_SESSION_CLOSE        = 609
	// Reserved up to 649.

)

type AuthSessionStart struct {
	Hostkey []byte
}

func (m AuthSessionStart) TypeId() uint16 {
	return AUTH_SESSION_START
}

func NewAuthSessionStart(data []byte) (AuthSessionStart, error) {

	m := AuthSessionStart{}
	m.Hostkey = make([]byte, len(data))
	copy(m.Hostkey, data)
	return m, nil
}

type AuthSessionHS1 struct {
	SessionId        uint32
	HandshakePayload []byte
}

func (m AuthSessionHS1) TypeId() uint16 {
	return AUTH_SESSION_HS1
}

func NewAuthSessionHS1(data []byte) (AuthSessionHS1, error) {

	m := AuthSessionHS1{}
	buf := bytes.NewReader(data)
	if err := binary.Read(buf, binary.BigEndian, &m.SessionId); err != nil {
		return m, err
	}
	data = data[4:]
	m.HandshakePayload = make([]byte, len(data))
	copy(m.HandshakePayload, data)
	return m, nil
}

type AuthSessionIncomingHS1 struct {
	Reserved         uint16
	HostkeyLength    uint16
	Hostkey          []byte
	HandshakePayload []byte
}

func (m AuthSessionIncomingHS1) TypeId() uint16 {
	return AUTH_SESSION_INCOMING_HS1
}

func NewAuthSessionIncomingHS1(data []byte) (AuthSessionIncomingHS1, error) {

	m := AuthSessionIncomingHS1{}
	reader := bytes.NewReader(data)
	if err := binary.Read(reader, binary.BigEndian, &m.Reserved); err != nil {
		return m, err
	}
	if err := binary.Read(reader, binary.BigEndian, &m.HostkeyLength); err != nil {
		return m, err
	}

	mustRead := int(m.HostkeyLength)
	hostkey := make([]byte, mustRead)
	if _, err := io.ReadFull(reader, hostkey); err != nil {
		return m, err
	}
	m.Hostkey = hostkey

	mustRead = len(data) - int(m.HostkeyLength) - 4
	payload := make([]byte, mustRead)
	if _, err := io.ReadFull(reader, payload); err != nil {
		return m, err
	}
	m.HandshakePayload = payload

	return m, nil
}

type AuthSessionHS2 struct {
	SessionId        uint32
	HandshakePayload []byte
}

func (m AuthSessionHS2) TypeId() uint16 {
	return AUTH_SESSION_HS2
}

func NewAuthSessionHS2(data []byte) (AuthSessionHS2, error) {

	var m AuthSessionHS2
	var reader *bytes.Reader
	var err error

	m = AuthSessionHS2{}
	reader = bytes.NewReader(data)
	if err = binary.Read(reader, binary.BigEndian, &m.SessionId); err != nil {
		return m, err
	}
	data = data[4:]
	m.HandshakePayload = make([]byte, len(data))
	copy(m.HandshakePayload, data)
	return m, nil
}

type AuthSessionIncomingHS2 struct {
	SessionId uint32
	Payload   []byte
}

func (m AuthSessionIncomingHS2) TypeId() uint16 {
	return AUTH_SESSION_INCOMING_HS2
}

func NewAuthSessionIncomingHS2(data []byte) (AuthSessionIncomingHS2, error) {

	var m AuthSessionIncomingHS2
	var reader *bytes.Reader
	var err error

	m = AuthSessionIncomingHS2{}
	reader = bytes.NewReader(data)

	if err = binary.Read(reader, binary.BigEndian, &m.SessionId); err != nil {
		return m, err
	}
	data = data[4:]

	m.Payload = make([]byte, len(data))
	copy(m.Payload, data)

	return m, nil
}

type AuthLayerEncrypt struct {
	Layers     uint8
	Reserved   uint8
	RequestId  uint16
	SessionIds []uint32
	Payload    []byte
}

func (m AuthLayerEncrypt) TypeId() uint16 {
	return AUTH_LAYER_ENCRYPT
}

func NewAuthLayerEncrypt(data []byte) (AuthLayerEncrypt, error) {

	var m AuthLayerEncrypt
	var reader *bytes.Reader
	var err error

	m = AuthLayerEncrypt{}
	reader = bytes.NewReader(data)

	if err = binary.Read(reader, binary.BigEndian, &m.Layers); err != nil {
		return m, err
	}
	if err = binary.Read(reader, binary.BigEndian, &m.Reserved); err != nil {
		return m, err
	}
	if err = binary.Read(reader, binary.BigEndian, &m.RequestId); err != nil {
		return m, err
	}

	m.SessionIds = make([]uint32, m.Layers)
	for i := 0; i < int(m.Layers); i++ {
		if err = binary.Read(reader, binary.BigEndian, &m.SessionIds[i]); err != nil {
			return m, err
		}
	}

	m.Payload = make([]byte, reader.Len())
	if _, err = io.ReadFull(reader, m.Payload); err != nil {
		return m, err
	}

	return m, nil
}

type AuthLayerEncryptResp struct {
	RequestId        uint16
	Reserved         uint16
	EncryptedPayload []byte
}

func (m AuthLayerEncryptResp) TypeId() uint16 {
	return AUTH_LAYER_ENCRYPT_RESP
}

func NewAuthLayerEncryptResp(data []byte) (AuthLayerEncryptResp, error) {

	var m AuthLayerEncryptResp
	var reader *bytes.Reader
	var err error

	m = AuthLayerEncryptResp{}
	reader = bytes.NewReader(data)

	if err = binary.Read(reader, binary.BigEndian, &m.RequestId); err != nil {
		return m, err
	}
	if err = binary.Read(reader, binary.BigEndian, &m.Reserved); err != nil {
		return m, err
	}

	m.EncryptedPayload = make([]byte, reader.Len())
	if _, err = io.ReadFull(reader, m.EncryptedPayload); err != nil {
		return m, err
	}

	return m, nil
}

type AuthLayerDecrypt struct {
	Layers           uint8
	Reserved         uint8
	RequestId        uint16
	SessionIds       []uint32
	EncryptedPayload []byte
}

func (m AuthLayerDecrypt) TypeId() uint16 {
	return AUTH_LAYER_DECRYPT
}

func NewAuthLayerDecrypt(data []byte) (AuthLayerDecrypt, error) {

	var m AuthLayerDecrypt
	var reader *bytes.Reader
	var err error

	m = AuthLayerDecrypt{}
	reader = bytes.NewReader(data)

	if err = binary.Read(reader, binary.BigEndian, &m.Layers); err != nil {
		return m, err
	}
	if err = binary.Read(reader, binary.BigEndian, &m.Reserved); err != nil {
		return m, err
	}
	if err = binary.Read(reader, binary.BigEndian, &m.RequestId); err != nil {
		return m, err
	}

	m.SessionIds = make([]uint32, m.Layers)
	for i := 0; i < int(m.Layers); i++ {
		if err = binary.Read(reader, binary.BigEndian, &m.SessionIds[i]); err != nil {
			return m, err
		}
	}

	m.EncryptedPayload = make([]byte, reader.Len())
	if _, err = io.ReadFull(reader, m.EncryptedPayload); err != nil {
		return m, err
	}

	return m, nil
}

type AuthLayerDecryptResp struct {
	RequestId        uint16
	Reserved         uint16
	DecryptedPayload []byte
}

func (m AuthLayerDecryptResp) TypeId() uint16 {
	return AUTH_LAYER_DECRYPT_RESP
}

func NewAuthLayerDecryptResp(data []byte) (AuthLayerDecryptResp, error) {

	var m AuthLayerDecryptResp
	var reader *bytes.Reader
	var err error

	m = AuthLayerDecryptResp{}
	reader = bytes.NewReader(data)

	if err = binary.Read(reader, binary.BigEndian, &m.RequestId); err != nil {
		return m, err
	}
	if err = binary.Read(reader, binary.BigEndian, &m.Reserved); err != nil {
		return m, err
	}

	m.DecryptedPayload = make([]byte, reader.Len())
	if _, err = io.ReadFull(reader, m.DecryptedPayload); err != nil {
		return m, err
	}

	return m, nil
}

type AuthSessionClose struct {
	SessionId uint32
}

func (m AuthSessionClose) TypeId() uint16 {
	return AUTH_SESSION_CLOSE
}

func NewAuthSessionClose(data []byte) (AuthSessionClose, error) {

	var m AuthSessionClose
	var reader *bytes.Reader
	var err error

	m = AuthSessionClose{}
	reader = bytes.NewReader(data)

	if err = binary.Read(reader, binary.BigEndian, &m.SessionId); err != nil {
		return m, err
	}
	return m, nil
}
