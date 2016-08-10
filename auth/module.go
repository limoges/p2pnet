package auth

import (
	"bytes"
	"crypto/rsa"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/limoges/p2pnet"
	"github.com/limoges/p2pnet/cfg"
	"github.com/limoges/p2pnet/msg"
)

// These constants are used for the module's configuration.
const (
	// The token identifying the module's configurations.
	ModuleToken = "ONION_AUTHENTICATION"
	// The token identifying the port to bind to.
	ApiAddrToken = "api_address"
	// The default API address to bind to.
	DefaultApiAddr = "127.0.0.1:7005"
	// The configurtion token providing the module with the hostkey file.
	HostkeyToken = "HOSTKEY"
	// The default location of the hostkey file.
	DefaultHostkey = "hostkey.pem"
)

var (
	ErrNoBlockFound = errors.New("No block found in key file")
)

// This module only communicates with the Onion module.
type Auth struct {
	PrivateKey *rsa.PrivateKey

	Sessions   map[uint32]*Session
	APIAddr    string
	ListenAddr string
}

func New(conf *cfg.Configurations) (*Auth, error) {

	var auth *Auth
	var priv *rsa.PrivateKey
	var err error
	var hostkeyPath string

	auth = &Auth{}
	conf.Init(&hostkeyPath, "", HostkeyToken, DefaultHostkey)
	conf.Init(&auth.APIAddr, ModuleToken, ApiAddrToken, DefaultApiAddr)

	if priv, err = ReadPEMPrivateKey(hostkeyPath); err != nil {
		fmt.Printf("Could not read necessary keys from '%v'.\n", hostkeyPath)
		return nil, err
	}
	auth.PrivateKey = priv
	auth.Sessions = make(map[uint32]*Session)
	return auth, nil
}

func (a *Auth) Name() string {
	return ModuleToken
}

func (a *Auth) Addresses() (APIAddr, P2PAddr string) {
	return a.APIAddr, ""
}

func (a *Auth) Run() error {

	select {
	// Run what should be ran immediately.
	// Don't return from this call until ready to end.
	}
}

func (a *Auth) Handle(source net.Conn, message msg.Message) error {

	switch message.(type) {
	case msg.AuthSessionStart:
		m := message.(msg.AuthSessionStart)
		return a.handleSessionStart(source, m)
	case msg.AuthSessionIncomingHS1:
		m := message.(msg.AuthSessionIncomingHS1)
		return a.handleSessionIncomingHS1(source, m)
	case msg.AuthSessionIncomingHS2:
		m := message.(msg.AuthSessionIncomingHS2)
		return a.handleSessionIncomingHS2(source, m)
	case msg.AuthSessionClose:
		m := message.(msg.AuthSessionClose)
		return a.handleSessionClose(source, m)
	case msg.AuthLayerEncrypt:
		m := message.(msg.AuthLayerEncrypt)
		return a.handleLayerEncrypt(source, m)
	case msg.AuthLayerDecrypt:
		m := message.(msg.AuthLayerDecrypt)
		return a.handleLayerDecrypt(source, m)
	default:
		return a.handleUnknown(source, message)
	}
	return nil
}

func (a *Auth) handleSessionStart(source net.Conn, m msg.AuthSessionStart) error {

	var sessionHS1 *msg.AuthSessionHS1
	var err error

	if sessionHS1, err = a.StartSession(m.Hostkey); err != nil {
		return err
	}
	return msg.Send(source, sessionHS1)
}

func (a *Auth) handleSessionIncomingHS1(source net.Conn, m msg.AuthSessionIncomingHS1) error {

	var sessionHS2 *msg.AuthSessionHS2
	var err error

	if sessionHS2, err = a.IncomingHandshake1(m.Hostkey, m.HandshakePayload); err != nil {
		return err
	}
	return msg.Send(source, sessionHS2)
}

func (a *Auth) handleSessionIncomingHS2(source net.Conn, m msg.AuthSessionIncomingHS2) error {

	var err error

	if err = a.IncomingHandshake2(m.SessionId, m.Payload); err != nil {
		return msg.Send(source, msg.AuthSessionDeclined{})
	}

	return msg.Send(source, msg.AuthSessionConfirmed{})
}

func (a *Auth) handleSessionClose(source net.Conn, m msg.AuthSessionClose) error {

	var id uint32
	id = m.SessionId
	a.CloseSession(id)
	return nil
}

func (a *Auth) handleLayerEncrypt(source net.Conn, m msg.AuthLayerEncrypt) error {

	var session *Session
	var present bool
	var payload []byte
	var encrypted []byte
	var err error

	payload = make([]byte, len(m.Payload))
	copy(payload, m.Payload)

	for _, sessionId := range m.SessionIds {
		if session, present = a.Sessions[sessionId]; !present {
			return errors.New(fmt.Sprintf("Session %v does not exist.", sessionId))
		}
		if encrypted, err = session.Encrypt(payload); err != nil {
			return errors.New("Could not encrypt payload")
		}
		payload = make([]byte, len(encrypted))
		copy(payload, encrypted)
	}

	var response msg.AuthLayerEncryptResp
	response = msg.AuthLayerEncryptResp{}
	response.RequestId = m.RequestId
	response.EncryptedPayload = make([]byte, len(payload))
	copy(response.EncryptedPayload, payload)

	return msg.Send(source, response)
}

func (a *Auth) handleLayerDecrypt(source net.Conn, m msg.AuthLayerDecrypt) error {

	var session *Session
	var present bool
	var payload []byte
	var decrypted []byte
	var err error

	payload = make([]byte, len(m.EncryptedPayload))
	copy(payload, m.EncryptedPayload)

	for i := len(m.SessionIds) - 1; i >= 0; i-- {
		sessionId := m.SessionIds[i]

		if session, present = a.Sessions[sessionId]; !present {
			return errors.New(fmt.Sprintf("Session %v does not exist.", sessionId))
		}
		if decrypted, err = session.Decrypt(payload); err != nil {
			return errors.New("Could not decrypt payload")
		}
		payload = make([]byte, len(decrypted))
		copy(payload, decrypted)
	}

	var response msg.AuthLayerDecryptResp
	response = msg.AuthLayerDecryptResp{}
	response.RequestId = m.RequestId
	response.DecryptedPayload = make([]byte, len(payload))
	copy(response.DecryptedPayload, payload)

	return msg.Send(source, response)
}

func (a *Auth) handleUnknown(source net.Conn, m msg.Message) error {

	fmt.Printf("Unhandled message: %v\n", msg.Identifier(m.TypeId()))
	return p2pnet.ErrModuleDoesNotHandle
}

func (a *Auth) StartSession(hostkey []byte) (*msg.AuthSessionHS1, error) {

	var pub *rsa.PublicKey
	var err error
	var session *Session
	var handshake1 *msg.AuthSessionHS1

	// First, check that the hostkey is a valid rsa.PublicKey
	if pub, err = ParsePublicKey(hostkey); err != nil {
		log.Println("Could not parse hostkey into public key format.")
		return nil, err
	}

	// Start creating the session.
	if session, err = NewSession(a); err != nil {
		log.Println("Could not create session.")
		return nil, err
	}

	if handshake1, err = session.CreateHandshake1(pub); err != nil {
		return nil, err
	}

	session.RemotePublicKey = pub
	a.Sessions[session.Id] = session
	return handshake1, nil
}

func (a *Auth) IncomingHandshake1(hostkey []byte, payload []byte) (*msg.AuthSessionHS2, error) {

	var pub *rsa.PublicKey
	var err error
	var handshake1 *msg.AuthHandshake1
	var session *Session
	var handshake2 *msg.AuthSessionHS2

	// Check that the remote hostkey is a valid rsa.PublicKey
	if pub, err = ParsePublicKey(hostkey); err != nil {
		log.Println("Could not parse hostkey into public key format.")
		return nil, err
	}

	// Parse the payload for the handshake message
	if handshake1, err = unloadHandshake1(payload); err != nil {
		log.Println("Could not parse handshake payload.")
		return nil, err
	}

	if session, err = NewIncomingSession(a, handshake1.EncryptedKey[:], handshake1.EncryptedHMAC[:]); err != nil {
		log.Println(err)
		return nil, err
	}

	if handshake2, err = session.CreateHandshake2(pub); err != nil {
		log.Println(err)
		return nil, err
	}

	session.RemotePublicKey = pub
	a.Sessions[session.Id] = session
	return handshake2, nil
}

func (a *Auth) IncomingHandshake2(id uint32, payload []byte) error {

	var session *Session
	var ok bool
	var err error
	var handshake2 *msg.AuthHandshake2

	// Check if the session exists
	if session, ok = a.Sessions[id]; !ok {
		return errors.New("Session does not exist")
	}

	// Validate the handshake payload
	if handshake2, err = unloadHandshake2(payload); err != nil {
		log.Println("Could not parse handshake payload.")
		return err
	}

	if err = session.DecryptRemoteHMAC(a.PrivateKey, handshake2.EncryptedHMAC[:]); err != nil {
		return err
	}
	return nil
}

func buildPayload(m msg.Message) ([]byte, error) {

	var buf *bytes.Buffer
	var err error

	buf = new(bytes.Buffer)
	if err = msg.Write(buf, m); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func unloadHandshake1(payload []byte) (*msg.AuthHandshake1, error) {

	var reader *bytes.Reader
	var message msg.Message
	var handshake1 msg.AuthHandshake1
	var err error
	var ok bool

	// First read the message sent through the payload.
	reader = bytes.NewReader(payload)

	if message, err = msg.Read(reader); err != nil {
		return nil, err
	}

	// We expect a AuthHandshake1 message in the payload.
	if handshake1, ok = message.(msg.AuthHandshake1); !ok {
		return nil, errors.New("Unexpected handshake payload")
	}

	return &handshake1, nil
}

func unloadHandshake2(payload []byte) (*msg.AuthHandshake2, error) {

	var reader *bytes.Reader
	var message msg.Message
	var handshake2 msg.AuthHandshake2
	var err error
	var ok bool

	// First read the message sent through the payload.
	reader = bytes.NewReader(payload)

	if message, err = msg.Read(reader); err != nil {
		return nil, err
	}

	// We expect a AuthHandshake1 message in the payload.
	if handshake2, ok = message.(msg.AuthHandshake2); !ok {
		return nil, errors.New("Unexpected handshake payload")
	}

	return &handshake2, nil
}

func (a *Auth) CloseSession(id uint32) {
	log.Println("CloseSession not implemented")
}
