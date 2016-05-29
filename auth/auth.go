package auth

import (
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/limoges/p2pnet/cfg"
	"github.com/limoges/p2pnet/msg"
	"io/ioutil"
	"math/rand"
)

// This module only communicates with the Onion module.
type Auth struct {
	HostkeyPath string
	block       *pem.Block

	Sessions map[uint32]string
}

var (
	ErrNoBlockFound = errors.New("No block found in key file")
)

// Reads a PEM-formatted RSA public key from the given path into
// a block format given by crypto/pem.
// If any error occurs, the block returned is set to nil.
func readPublicKeyFromFile(path string) (block *pem.Block, err error) {

	data, err := ioutil.ReadFile(auth.HostkeyPath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, ErrNoBlockFound
	}

	return block, nil
}

func New(conf *cfg.Configurations) (auth *Auth, err error) {

	auth = &Auth{}
	conf.Init(&auth.HostkeyPath, "", "HOSTKEY", "hostkey.pem")

	block, err := readPublicKeyFromFile(auth.HostkeyPath)
	if err != nil {
		return nil, err
	}

	auth.block = block
	return auth, nil
}

func (a *Auth) Run() error {

	apiListener, err := net.Listen("tcp", a.apiAddr)
	if err != nil {
		fmt.Println(err)
		return err
	}

	a.listenAPI(apiListener)
}

func (a *Auth) ListenAPI(ln net.Listener) {

	fmt.Printf("Launched API Listener on %v\n", ln.Addr())
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

	// switch t := m.(type) {
	// case msg.AuthSessionStart:
	// 	am := msg.AuthSessionStart(m)
	// 	a.StartSessionKeyEstablishment(am.Hostkey)
	// default:
	// 	fmt.Printf("message is not handled by Auth:%v\n", message)
	// }
}

func (a *Auth) createSessionID() uint32 {
	return rand.Uint32()
}

func (a *Auth) SessionStart(hostkey []byte) {

	sessionID := a.createSessionID()

	// Create handshake payload
	a.Sessions[sessionID] = string(hostkey)

	m := AuthSessionHS1{
		SessionID:        sessionID,
		HandshakePayload: {},
	}
}
