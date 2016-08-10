package msg

import (
	"bytes"
	"io"
)

const (
	AUTH_HANDSHAKE1        = 700
	AUTH_HANDSHAKE2        = 701
	AUTH_SESSION_CONFIRMED = 702
	AUTH_SESSION_DECLINED  = 703
)

type AuthHandshake1 struct {
	EncryptedKey  [512]byte
	EncryptedHMAC [512]byte
}

func (m AuthHandshake1) TypeId() uint16 {
	return AUTH_HANDSHAKE1
}

func NewAuthHandshake1(data []byte) (AuthHandshake1, error) {

	var m AuthHandshake1
	var reader *bytes.Reader
	var err error

	m = AuthHandshake1{}
	reader = bytes.NewReader(data)

	if _, err = io.ReadFull(reader, m.EncryptedKey[:]); err != nil {
		return m, err
	}
	if _, err = io.ReadFull(reader, m.EncryptedHMAC[:]); err != nil {
		return m, err
	}
	return m, nil
}

type AuthHandshake2 struct {
	EncryptedHMAC [512]byte
}

func (m AuthHandshake2) TypeId() uint16 {
	return AUTH_HANDSHAKE2
}

func NewAuthHandshake2(data []byte) (AuthHandshake2, error) {

	var m AuthHandshake2
	var reader *bytes.Reader
	var err error

	m = AuthHandshake2{}
	reader = bytes.NewReader(data)

	if _, err = io.ReadFull(reader, m.EncryptedHMAC[:]); err != nil {
		return m, err
	}
	return m, nil
}

type AuthSessionConfirmed struct {
}

func (m AuthSessionConfirmed) TypeId() uint16 {
	return AUTH_SESSION_CONFIRMED
}

func NewAuthSessionConfirmed(data []byte) (AuthSessionConfirmed, error) {
	return AuthSessionConfirmed{}, nil
}

type AuthSessionDeclined struct {
}

func (m AuthSessionDeclined) TypeId() uint16 {
	return AUTH_SESSION_CONFIRMED
}

func NewAuthSessionDeclined(data []byte) (AuthSessionDeclined, error) {
	return AuthSessionDeclined{}, nil
}
