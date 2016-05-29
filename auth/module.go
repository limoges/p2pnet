package auth

import (
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"time"

	"github.com/limoges/p2pnet"
	"github.com/limoges/p2pnet/cfg"
	"github.com/limoges/p2pnet/msg"
)

const (
	moduleToken    = "ONION_AUTHENTICATION"
	apiAddrToken   = "api_address"
	defaultApiAddr = "127.0.0.1:7005"
	hostkeyToken   = "HOSTKEY"
	defaultHostkey = "hostkey.pem"
)

var (
	ErrNoBlockFound = errors.New("No block found in key file")
)

// This module only communicates with the Onion module.
type Auth struct {
	HostkeyPath string
	Hostkey     []byte

	Sessions   map[p2pnet.SessionId]p2pnet.Session
	APIAddr    string
	ListenAddr string
}

// Reads a PEM-formatted RSA public key from the given path into
// a block format given by crypto/pem.
// If any error occurs, the block returned is set to nil.
func readPublicKeyFromFile(path string) (block *pem.Block, err error) {

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ = pem.Decode(data)
	if block == nil {
		return nil, ErrNoBlockFound
	}

	return block, nil
}

func New(conf *cfg.Configurations) (auth *Auth, err error) {

	auth = &Auth{}
	conf.Init(&auth.HostkeyPath, "", hostkeyToken, defaultHostkey)
	conf.Init(&auth.APIAddr, moduleToken, apiAddrToken, defaultApiAddr)

	block, err := readPublicKeyFromFile(auth.HostkeyPath)
	if err != nil {
		return nil, err
	}

	auth.Hostkey = make([]byte, len(block.Bytes))
	copy(auth.Hostkey, block.Bytes)
	fmt.Printf("The DER key's length in bytes is %v\n", len(auth.Hostkey))
	auth.Sessions = make(map[p2pnet.SessionId]p2pnet.Session)
	return auth, nil
}

func (a *Auth) Name() string {
	return moduleToken
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
	case msg.AuthSessionStart, *msg.AuthSessionStart:
		m := message.(msg.AuthSessionStart)
		a.StartSession(source, m.Hostkey)
	case msg.AuthSessionHS1, *msg.AuthSessionHS1:
		m := message.(msg.AuthSessionHS1)
		fmt.Println(m)
	case msg.AuthSessionIncomingHS1, *msg.AuthSessionIncomingHS1:
		m := message.(msg.AuthSessionIncomingHS1)
		fmt.Println(m)
	default:
		fmt.Printf("Unhandled message: %v\n", message)
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

	// We generate a session id which is not currently being used.
	sessionID := a.generateUnusedSessionID()

	// Create the peer, holding his address and hostkey.
	peer := p2pnet.Peer{
		Hostkey: hostkey,
		Addr:    conn.RemoteAddr().String(),
	}

	// We associate the peer with a session
	session := p2pnet.Session{
		ID:      sessionID,
		Partner: peer,
	}

	// Save the session into our cache
	a.Sessions[sessionID] = session

	// Send the response to the session start.
	m := msg.AuthSessionHS1{
		SessionId:        uint32(sessionID),
		HandshakePayload: []byte{},
	}

	msg.WriteMessage(conn, m)

	return nil
}

func (a *Auth) CloseSession(id p2pnet.SessionId) {
	delete(a.Sessions, id)
}
