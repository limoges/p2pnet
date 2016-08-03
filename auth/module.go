package auth

import (
	"bytes"
	"crypto/rsa"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/limoges/p2pnet"
	"github.com/limoges/p2pnet/cfg"
	"github.com/limoges/p2pnet/msg"
)

// These constants are used for the module's configuration.
const (
	// The token identifying the module's configurations.
	ModuleToken = "ONION_AUTHENTICATION"
	// The token identifying the port to bind to.
	ApiAddrToken   = "api_address"
	DefaultApiAddr = "127.0.0.1:7005"
	HostkeyToken   = "HOSTKEY"
	DefaultHostkey = "hostkey.pem"
)

var (
	ErrNoBlockFound = errors.New("No block found in key file")
)

type Session struct {
	Key []byte
}

// This module only communicates with the Onion module.
type Auth struct {
	// HostkeyPath should refer to a file which contains a RSA private key.
	HostkeyPath string
	Keys        *Encryption

	Sessions   map[p2pnet.SessionId]Session
	APIAddr    string
	ListenAddr string
}

func New(conf *cfg.Configurations) (auth *Auth, err error) {

	auth = &Auth{}
	conf.Init(&auth.HostkeyPath, "", HostkeyToken, DefaultHostkey)
	conf.Init(&auth.APIAddr, ModuleToken, ApiAddrToken, DefaultApiAddr)

	keys, err := ReadKeys(auth.HostkeyPath)
	if err != nil {
		fmt.Printf("Could not read necessary keys from '%v'.\n", auth.HostkeyPath)
		return nil, err
	}
	auth.Keys = keys
	auth.Sessions = make(map[p2pnet.SessionId]Session)
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
		// Upon reception of a AUTH_SESSION_START from the Onion,
		// we generate the Handshake1.
		m := message.(msg.AuthSessionStart)
		a.StartSession(source, m.Hostkey)
	case msg.AuthSessionIncomingHS1:
		m := message.(msg.AuthSessionIncomingHS1)
		a.RespondIncomingHS1(source, m.Hostkey, m.HandshakePayload)
	case msg.AuthSessionIncomingHS2:
		fmt.Println("AuthSessionIncomingHS2")
		m := message.(msg.AuthSessionIncomingHS2)
		//a.RespondIncomingHS2(source, m)
		fmt.Println(m.TypeId())
	default:
		fmt.Printf("Unhandled message: %v\n", msg.Identifier(message.TypeId()))
		return p2pnet.ErrModuleDoesNotHandle
	}

	return nil
}

func (a *Auth) generateUnusedSessionID() p2pnet.SessionId {

	// With the estimated number of simultaneous connexion very low,
	// it is suggested to use a 32 bits session ID which will have
	// about odds of a collision in about 1 in 10 millions for 30 or
	// simultaneous sessions.
	// Therefore, using the time as a random number generator seed seems
	// perfectly fine.

	seed := time.Now().UnixNano()
	source := rand.NewSource(seed)
	rng := rand.New(source)

	// We generate a random number and then try again until we can find
	// a number that is not currently in use. Again, there shouldn't be too
	// many collisions locally.
	sessionID := p2pnet.SessionId(rng.Uint32())
	count := 0
	for {
		if _, alreadyInUse := a.Sessions[sessionID]; !alreadyInUse {
			break
		}
		count = count + 1
		sessionID = p2pnet.SessionId(rng.Uint32())
	}
	return sessionID
}

func (a *Auth) StartSession(conn net.Conn, hostkey []byte) error {

	var sessionId p2pnet.SessionId
	var hs1 msg.AuthSessionHS1
	var pub *rsa.PublicKey
	var err error
	var key []byte
	var ciphertext []byte
	var session Session
	var buf *bytes.Buffer

	// We generate a session id which is not currently being used.
	sessionId = a.generateUnusedSessionID()

	// Parse the given hostkey into a RSA Public Key
	if pub, err = ParsePublicKey(hostkey); err != nil {
		return err
	}

	// Generate a symmetric key to share. This is our payload.
	if key, err = GenerateNewSymmetricKey(); err != nil {
		return err
	}

	// Encrypt the given key using asymmetric encryption.
	if ciphertext, err = AsymmetricEncrypt(pub, key); err != nil {
		return err
	}

	buf = new(bytes.Buffer)
	handshake1 := msg.AuthHandshake1{
		Cipher: ciphertext,
	}

	if err = msg.WriteMessage(buf, handshake1); err != nil {
		return err
	}

	// Send the response to the session start.
	hs1 = msg.AuthSessionHS1{
		SessionId:        uint32(sessionId),
		HandshakePayload: buf.Bytes(),
	}

	// Store the session's key.
	session = Session{
		Key: key,
	}

	// Associate the session ID with a peer's hostkey
	a.Sessions[sessionId] = session

	log.Println("Sending session handshake1 to onion...")
	msg.WriteMessage(conn, hs1)

	return nil
}

func (a *Auth) RespondIncomingHS1(conn net.Conn, hostkey, payload []byte) error {

	var key []byte
	var err error
	var message msg.Message
	var hs1 msg.AuthHandshake1
	var handshake2 msg.AuthHandshake2
	var hs2 msg.AuthSessionHS2
	var reader io.Reader
	var session Session
	var ok bool
	var sessionId p2pnet.SessionId
	var pub *rsa.PublicKey
	var ciphertext []byte
	var buf *bytes.Buffer

	// First read the message sent through the payload.
	reader = bytes.NewReader(payload)

	if message, err = msg.ReadMessage(reader); err != nil {
		return err
	}

	// We expect a AuthHandshake1 message in the payload.
	if hs1, ok = message.(msg.AuthHandshake1); !ok {
		return errors.New("Unexpected handshake payload")
	}

	// We then try to get the session key presented by the peer, using our
	// private key to decrypt the cipher containing the session key.
	if key, err = AsymmetricDecrypt(a.Keys.PrivateKey, hs1.Cipher); err != nil {
		return err
	}

	// We then generate our own internal session ID which we'll use to refer to
	// this session.
	sessionId = a.generateUnusedSessionID()

	session = Session{
		Key: key,
	}

	// Store the key into it's associated session.
	a.Sessions[sessionId] = session

	if pub, err = ParsePublicKey(hostkey); err != nil {
		return err
	}

	if ciphertext, err = AsymmetricEncrypt(pub, key); err != nil {
		return err
	}

	handshake2 = msg.AuthHandshake2{
		Cipher: ciphertext,
	}

	buf = new(bytes.Buffer)

	if err = msg.WriteMessage(buf, handshake2); err != nil {
		return err
	}

	hs2 = msg.AuthSessionHS2{
		SessionId:        uint32(sessionId),
		HandshakePayload: buf.Bytes(),
	}

	// Encode the response message
	if err = msg.WriteMessage(conn, hs2); err != nil {
		return err
	}

	return nil
}

func (a *Auth) CloseSession(id p2pnet.SessionId) {
	delete(a.Sessions, id)
}
