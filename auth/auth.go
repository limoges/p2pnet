package auth

import (
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/limoges/p2pnet/cfg"
	"github.com/limoges/p2pnet/msg"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"time"
)

type SessionId uint32

// This module only communicates with the Onion module.
type Auth struct {
	HostkeyPath string
	block       *pem.Block

	Sessions   map[SessionId]Session
	APIAddr    string
	ListenAddr string
}

type Session struct {
	ID      SessionId
	Partner Peer
}

type Peer struct {
	Hostkey []byte
	Addr    string
}

var (
	ErrNoBlockFound = errors.New("No block found in key file")
)

const (
	moduleToken  = "ONION_AUTHENTICATION"
	apiAddrToken = "api_address"
	//listenAddrToken = "listen_addr"
)

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
	conf.Init(&auth.HostkeyPath, "", "HOSTKEY", "hostkey.pem")
	conf.Init(&auth.APIAddr, moduleToken, apiAddrToken, "127.0.0.1:7005")

	block, err := readPublicKeyFromFile(auth.HostkeyPath)
	if err != nil {
		return nil, err
	}

	auth.block = block
	auth.Sessions = make(map[SessionId]Session)
	return auth, nil
}

func (a *Auth) Run() error {

	apiListener, err := net.Listen("tcp", a.APIAddr)
	if err != nil {
		fmt.Println(err)
		return err
	}

	listener, err := net.Listen("tcp", a.ListenAddr)
	if err != nil {
		fmt.Println(err)
		return err
	}

	go a.listen(apiListener)
	fmt.Printf("Launched API Listener on %v\n", apiListener.Addr())
	go a.listen(listener)
	fmt.Printf("Launched Listener on %v\n", listener.Addr())

	for {
		// Run what should be ran immediately.
		// Don't return from this call until ready to end.
	}
}

func (a *Auth) listen(ln net.Listener) {

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
		} else {
			go a.handle(conn)
		}
	}
}

func (a *Auth) handle(conn net.Conn) {

	fmt.Printf("New connexion from %v\n", conn.RemoteAddr())
	for {
		message, err := msg.ReadMessage(conn)
		if err != nil {
			if err == io.EOF {
				return
			}
			fmt.Println(err)
			return
		}

		switch message.(type) {
		case msg.AuthSessionStart, *msg.AuthSessionStart:
			m := message.(msg.AuthSessionStart)
			a.StartSession(conn, m.Hostkey)
		case msg.AuthSessionHS1, *msg.AuthSessionHS1:
			m := message.(msg.AuthSessionHS1)
			fmt.Println(m)
		case msg.AuthSessionIncomingHS1, *msg.AuthSessionIncomingHS1:
			m := message.(msg.AuthSessionIncomingHS1)
			fmt.Println(m)
		default:
			fmt.Printf("Unhandled message: %v\n", message)
		}
	}
}

func (a *Auth) generateUnusedSessionID() SessionId {

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
	sessionID := SessionId(rng.Uint32())
	for {
		if _, alreadyInUse := a.Sessions[sessionID]; !alreadyInUse {
			break
		}
		sessionID = SessionId(rng.Uint32())
	}
	return sessionID
}

func (a *Auth) StartSession(conn net.Conn, hostkey []byte) {

	// We generate a session id which is not currently being used.
	sessionID := a.generateUnusedSessionID()

	// Create the peer, holding his address and hostkey.
	peer := Peer{
		Hostkey: hostkey,
		Addr:    conn.RemoteAddr().String(),
	}

	// We associate the peer with a session
	session := Session{
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
}

func (a *Auth) CloseSession(id SessionId) {
	delete(a.Sessions, id)
}
