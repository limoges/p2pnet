package onion

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"strconv"

	"github.com/limoges/p2pnet"
	"github.com/limoges/p2pnet/auth"
	"github.com/limoges/p2pnet/cfg"
	"github.com/limoges/p2pnet/msg"
)

const (
	ModuleToken       = "ONION_FORWARDING"
	MinHopToken       = "min_hop_count"
	DefaultMinHop     = 2
	HopCountToken     = "hop_count"
	DefaultHopCount   = 5
	ListenAddrToken   = "listen_address"
	DefaultListenAddr = "127.0.0.1:7014"
	ApiAddrToken      = "api_address"
	DefaultApiAddr    = "127.0.0.1:7004"
	HostkeyToken      = "HOSTKEY"
	DefaultHostkey    = "hostkey.pem"
)

type Onion struct {
	MinimalHopCount int
	HopCount        int
	HostkeyPath     string
	Hostkey         []byte

	ListenAddr string
	APIAddr    string
	AuthAddr   string

	Peers    map[p2pnet.Identity]string
	Sessions map[p2pnet.SessionId]p2pnet.Identity
}

func New(conf *cfg.Configurations) (*Onion, error) {

	var mod *Onion
	var keys *auth.Encryption
	var err error

	mod = &Onion{}
	conf.Init(&mod.MinimalHopCount, ModuleToken, MinHopToken, DefaultMinHop)
	conf.Init(&mod.HopCount, ModuleToken, HopCountToken, DefaultHopCount)
	conf.Init(&mod.ListenAddr, ModuleToken, ListenAddrToken, DefaultListenAddr)
	conf.Init(&mod.APIAddr, ModuleToken, ApiAddrToken, DefaultApiAddr)
	conf.Init(&mod.AuthAddr, auth.ModuleToken, auth.ApiAddrToken, auth.DefaultApiAddr)
	conf.Init(&mod.HostkeyPath, "", HostkeyToken, DefaultHostkey)

	if keys, err = auth.ReadKeys(mod.HostkeyPath); err != nil {
		return nil, err
	}
	mod.Hostkey = keys.Hostkey
	mod.Peers = make(map[p2pnet.Identity]string)
	mod.Sessions = make(map[p2pnet.SessionId]p2pnet.Identity)
	return mod, nil
}

func (o *Onion) Name() string {
	return ModuleToken
}

func (o *Onion) Addresses() (APIAddr, P2PAddr string) {
	return o.APIAddr, o.ListenAddr
}

func (o *Onion) Run() error {
	select {}
}

func (o *Onion) Handle(source net.Conn, message msg.Message) error {

	switch message.(type) {
	case msg.OnionTunnelBuild:
		m := message.(msg.OnionTunnelBuild)
		return o.BuildTunnel(source, m.IPAddr, int(m.Port), m.DstHostkey)
	case msg.AuthSessionIncomingHS1:
		m := message.(msg.AuthSessionIncomingHS1)
		return o.RespondHandshake1(source, m)
	default:
		fmt.Printf("Unhandled message type:%v\n", msg.Identifier(message.TypeId()))
		return p2pnet.ErrModuleDoesNotHandle
	}
}

func (o *Onion) storeIdentity(hostkey []byte, hostport string) {

	var identity p2pnet.Identity

	// Get the hostkey's identity (sha256 sum of the hostkey) as a string
	identity = p2pnet.GetIdentity(hostkey)

	o.Peers[identity] = hostport
}

func (o *Onion) storeSessionId(hostkey []byte, id p2pnet.SessionId) {

	var identity p2pnet.Identity

	identity = p2pnet.GetIdentity(hostkey)

	o.Sessions[id] = identity
}

func (o *Onion) BuildTunnel(source net.Conn, addr net.IP, port int, hostkey []byte) error {

	var handshake1 *msg.AuthSessionHS1
	var response *msg.AuthHandshake2
	var incomingHS2 msg.AuthSessionIncomingHS2
	var hostport string
	var err error
	var sessionId p2pnet.SessionId
	var buf *bytes.Buffer

	// Convert the IP address and port to a string
	hostport = net.JoinHostPort(addr.String(), strconv.Itoa(port))

	o.storeIdentity(hostkey, hostport)

	if handshake1, err = o.requestSessionStart(hostkey); err != nil {
		return err
	}

	sessionId = p2pnet.SessionId(handshake1.SessionId)
	o.storeSessionId(hostkey, sessionId)

	if response, err = o.requestDestinationHandshake(hostport, *handshake1); err != nil {
		return err
	}

	buf = new(bytes.Buffer)
	if err = msg.WriteMessage(buf, response); err != nil {
		return err
	}

	incomingHS2 = msg.AuthSessionIncomingHS2{
		SessionId: uint32(sessionId),
		Payload:   buf.Bytes(),
	}

	if _, err = o.forwardAuth(incomingHS2); err != nil {
		return err
	}

	return nil
}

func (o *Onion) requestDestinationHandshake(hostport string, hs1 msg.AuthSessionHS1) (*msg.AuthHandshake2, error) {

	var conn net.Conn
	var err error
	var request msg.AuthSessionIncomingHS1
	var response msg.Message
	var hs2 msg.AuthHandshake2
	var ok bool

	if conn, err = net.Dial("tcp", hostport); err != nil {
		return nil, err
	}
	defer conn.Close()

	request = msg.AuthSessionIncomingHS1{
		HostkeyLength:    uint16(len(o.Hostkey)),
		Hostkey:          o.Hostkey,
		HandshakePayload: hs1.HandshakePayload,
	}

	if err = msg.WriteMessage(conn, request); err != nil {
		return nil, err
	}

	if response, err = msg.ReadMessage(conn); err != nil {
		return nil, err
	}

	if hs2, ok = response.(msg.AuthHandshake2); !ok {
		return nil, errors.New("Unexpected message type")
	}

	return &hs2, nil
}

func (o *Onion) requestSessionStart(hostkey []byte) (*msg.AuthSessionHS1, error) {

	var conn net.Conn
	var err error
	var request msg.AuthSessionStart
	var response msg.Message
	var sessionHS1 msg.AuthSessionHS1
	var ok bool

	if conn, err = net.Dial("tcp", o.AuthAddr); err != nil {
		return nil, err
	}
	defer conn.Close()

	request = msg.AuthSessionStart{
		Hostkey: hostkey,
	}

	if err = msg.WriteMessage(conn, request); err != nil {
		return nil, err
	}

	if response, err = msg.ReadMessage(conn); err != nil {
		return nil, err
	}

	if sessionHS1, ok = response.(msg.AuthSessionHS1); !ok {
		return nil, errors.New("Expecting AuthSessionHS1 response")
	}

	return &sessionHS1, nil
}

func (o *Onion) forwardAuth(message msg.Message) (msg.Message, error) {

	var conn net.Conn
	var err error

	if conn, err = net.Dial("tcp", o.AuthAddr); err != nil {
		return nil, err
	}
	defer conn.Close()

	if err = msg.WriteMessage(conn, message); err != nil {
		return nil, err
	}

	return msg.ReadMessage(conn)
}

func (o *Onion) forwardAuthPayload(payload []byte) (msg.Message, error) {

	var reader *bytes.Reader
	var err error
	var message msg.Message

	reader = bytes.NewReader(payload)

	if message, err = msg.ReadMessage(reader); err != nil {
		return nil, err
	}

	return o.forwardAuth(message)
}

func (o *Onion) RespondHandshake1(conn net.Conn, message msg.AuthSessionIncomingHS1) error {

	var err error
	var request msg.Message
	var hs2 msg.AuthSessionHS2
	var response msg.AuthSessionIncomingHS2
	var ok bool
	var hostkey []byte

	if request, err = o.forwardAuth(message); err != nil {
		return err
	}

	if hs2, ok = request.(msg.AuthSessionHS2); !ok {
		fmt.Println(request.TypeId())
		return errors.New("Expecting message AuthSessionHS2")
	}

	response = msg.AuthSessionIncomingHS2{}
	if err = msg.WriteMessage(conn, response); err != nil {
		return err
	}

	hostkey = message.Hostkey

	o.storeSessionId(hostkey, p2pnet.SessionId(hs2.SessionId))

	return nil
}
