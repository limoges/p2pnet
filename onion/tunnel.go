package onion

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/limoges/p2pnet"
	"github.com/limoges/p2pnet/msg"
)

type Tunnel struct {
	Id    uint32
	onion *Onion
}

func NewTunnel(o *Onion) (*Tunnel, error) {

	var tunnel *Tunnel
	var err error

	tunnel = &Tunnel{}
	if tunnel.Id, err = o.localUnusedTunnelId(); err != nil {
		return nil, err
	}
	tunnel.onion = o

	return tunnel, nil
}

func (t *Tunnel) AddLink(hostport string, hostkey []byte) error {

	var err error

	t.onion.storeIdentity(hostport, hostkey)

	if err = t.onion.buildSession(hostport, hostkey); err != nil {
		return err
	}

	return nil
}

func (t *Tunnel) CreateTunnelReady() (*msg.OnionTunnelReady, error) {

	var tunnelReady *msg.OnionTunnelReady

	tunnelReady = &msg.OnionTunnelReady{
		TunnelId: t.Id,
	}
	return tunnelReady, nil
}

func (o *Onion) localUnusedTunnelId() (uint32, error) {

	const MaximumTunnelIdAttempts = 100
	var id uint32
	var alreadyInUse bool
	var count int

	seed := time.Now().UnixNano()
	source := rand.NewSource(seed)
	rng := rand.New(source)

	id = rng.Uint32()
	count = 1
	for {
		if count == MaximumTunnelIdAttempts {
			return 0, errors.New("Could not generate a new tunnel id")
		}
		if _, alreadyInUse = o.Tunnels[id]; !alreadyInUse {
			break
		}
		id = rng.Uint32()
		count = count + 1
	}
	return id, nil
}

func (o *Onion) repackageHandshake1(m *msg.AuthSessionHS1) *msg.AuthSessionIncomingHS1 {

	var repackaged *msg.AuthSessionIncomingHS1

	repackaged = &msg.AuthSessionIncomingHS1{}
	repackaged.HostkeyLength = uint16(len(o.Hostkey))
	repackaged.Hostkey = make([]byte, len(o.Hostkey))
	copy(repackaged.Hostkey, o.Hostkey)
	repackaged.HandshakePayload = make([]byte, len(m.HandshakePayload))
	copy(repackaged.HandshakePayload, m.HandshakePayload)

	return repackaged
}

func (o *Onion) repackageHandshake2(sessionId uint32, m *msg.AuthHandshake2) (*msg.AuthSessionIncomingHS2, error) {

	var buf *bytes.Buffer
	var repackaged msg.AuthSessionIncomingHS2
	var err error

	buf = new(bytes.Buffer)

	if err = msg.Write(buf, m); err != nil {
		return nil, err
	}

	repackaged = msg.AuthSessionIncomingHS2{}
	repackaged.SessionId = sessionId
	repackaged.Payload = make([]byte, buf.Len())
	copy(repackaged.Payload, buf.Bytes())

	return &repackaged, nil
}

func (o *Onion) packageIncomingHandshake2(sessionId uint32, m msg.Message) (*msg.AuthSessionIncomingHS2, error) {

	var buf *bytes.Buffer
	var incoming *msg.AuthSessionIncomingHS2
	var err error

	buf = new(bytes.Buffer)
	if err = msg.Write(buf, m); err != nil {
		return nil, err
	}

	incoming = &msg.AuthSessionIncomingHS2{
		SessionId: uint32(sessionId),
		Payload:   buf.Bytes(),
	}

	return incoming, nil
}
func (o *Onion) BuildTunnel(hopCount int, hostport string, hostkey []byte) (*msg.OnionTunnelReady, error) {

	var tunnelReady *msg.OnionTunnelReady
	var err error
	var tunnel *Tunnel

	if tunnel, err = NewTunnel(o); err != nil {
		return nil, err
	}

	// for _, i := range hopCount {
	// 		Building minim hops
	//		if err = tunnel.AddLink(randomHostport, randomHostkey); err != nil {
	//			return err
	//		}
	// }

	// Add the final link
	if err = tunnel.AddLink(hostport, hostkey); err != nil {
		fmt.Println(err)
		return nil, err
	}

	if tunnelReady, err = tunnel.CreateTunnelReady(); err != nil {
		return nil, err
	}

	return tunnelReady, nil
}

func (o *Onion) requestHandshake1(hostkey []byte) (*msg.AuthSessionHS1, error) {

	var request msg.AuthSessionStart
	var response msg.Message
	var validResponse msg.AuthSessionHS1
	var valid bool
	var err error

	// Build the Session Start request
	request = msg.AuthSessionStart{}
	request.Hostkey = make([]byte, len(hostkey))
	copy(request.Hostkey, hostkey)

	// Request the session start from the Auth module
	if response, err = requestFrom(o.AuthAddr, request); err != nil {
		return nil, err
	}

	// Validate the session start
	if validResponse, valid = response.(msg.AuthSessionHS1); !valid {
		return nil, errors.New("Invalid response to Auth Session Start")
	}

	return &validResponse, nil
}

func (o *Onion) requestHandshake2(request *msg.AuthSessionIncomingHS1) (*msg.AuthSessionHS2, error) {

	var response msg.Message
	var validResponse msg.AuthSessionHS2
	var valid bool
	var err error

	// Request the session start from the Auth module
	if response, err = requestFrom(o.AuthAddr, request); err != nil {
		return nil, err
	}

	// Validate the session start
	if validResponse, valid = response.(msg.AuthSessionHS2); !valid {
		return nil, errors.New("Invalid response to Auth Session Start")
	}

	return &validResponse, nil
}

func (o *Onion) finalHandshake(hostport string, handshake1 *msg.AuthSessionIncomingHS1) (*msg.AuthHandshake2, error) {

	var response msg.Message
	var validResponse msg.AuthHandshake2
	var valid bool
	var err error

	if response, err = requestFrom(hostport, handshake1); err != nil {
		return nil, err
	}

	if validResponse, valid = response.(msg.AuthHandshake2); !valid {
		return nil, errors.New("Invalid response expected AuthHandshake2")
	}

	return &validResponse, nil
}

func (o *Onion) buildSession(hostport string, hostkey []byte) error {

	var handshake1 *msg.AuthSessionHS1
	var repackaged1 *msg.AuthSessionIncomingHS1
	var handshake2 *msg.AuthHandshake2
	var repackaged2 *msg.AuthSessionIncomingHS2
	var response msg.Message
	var sessionId uint32
	var err error

	// Start the session with our Auth module
	if handshake1, err = o.requestHandshake1(hostkey); err != nil {
		return err
	}

	// Save the session Id and hostport.
	sessionId = handshake1.SessionId
	o.storeSession(sessionId, hostkey)
	o.storeIdentity(hostport, hostkey)

	// Repackage the handshake to send it.
	repackaged1 = o.repackageHandshake1(handshake1)

	// Request the handshake2 from the remote.
	if handshake2, err = o.finalHandshake(hostport, repackaged1); err != nil {
		return err
	}

	// Repackage the response and send it to auth.
	if repackaged2, err = o.repackageHandshake2(sessionId, handshake2); err != nil {
		return err
	}

	// Send the final handshake to the Auth module.
	if response, err = requestFrom(o.AuthAddr, repackaged2); err != nil {
		return err
	}

	// Check the message type session confirmed.
	if response.TypeId() != msg.AUTH_SESSION_CONFIRMED {
		return errors.New("Session has been denied.")
	}

	return nil
}

// Send a message and waits for the response.
func requestFrom(hostport string, message msg.Message) (msg.Message, error) {

	var conn net.Conn
	var err error

	// Attempting to connect with the host
	if conn, err = net.Dial("tcp", hostport); err != nil {
		return nil, err
	}

	// Close the connexion upon leaving the function.
	defer conn.Close()

	return msg.SendReceive(conn, message)
}

// Sends a message but does not wait for the response.
func forwardTo(hostport string, message msg.Message) error {

	var conn net.Conn
	var err error

	// Attempting to connect with the host
	if conn, err = net.Dial("tcp", hostport); err != nil {
		return err
	}

	// Close the connexion upon leaving the function.
	defer conn.Close()

	return msg.Send(conn, message)
}

func (o *Onion) storeIdentity(hostport string, hostkey []byte) {

	var identity p2pnet.Identity
	identity = p2pnet.GetIdentity(hostkey)
	o.Peers[identity] = hostport
}

func (o *Onion) storeSession(id uint32, hostkey []byte) {

	var identity p2pnet.Identity
	identity = p2pnet.GetIdentity(hostkey)
	o.Sessions[id] = identity
}
