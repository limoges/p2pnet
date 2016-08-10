package onion

import (
	"crypto/rsa"
	"fmt"
	"net"
	"strconv"

	"github.com/limoges/p2pnet"
	"github.com/limoges/p2pnet/auth"
	"github.com/limoges/p2pnet/cfg"
	"github.com/limoges/p2pnet/msg"
)

// These constants are used for the module's configurations.
const (
	// The token identifying the module's configurations.
	ModuleToken = "ONION_FORWARDING"
	// The token identifying the minimal hop count configuration.
	MinHopToken = "min_hop_count"
	// The default minimal hop count configuration.
	DefaultMinHop = 2
	// The token identifying the hop count configuration.
	HopCountToken = "hop_count"
	// The default hop count configuration.
	DefaultHopCount = 5
	// The token identifying the listen address configuration.
	ListenAddrToken = "listen_address"
	// The default listen address configuration.
	DefaultListenAddr = "127.0.0.1:7014"
	// The token identifying the api address configuration.
	ApiAddrToken = "api_address"
	// The default api address configuration.
	DefaultApiAddr = "127.0.0.1:7004"
	// The token identifying the hostkey configuration.
	HostkeyToken = "HOSTKEY"
	// The default hostkey configuration.
	DefaultHostkey = "hostkey.pem"
)

type Onion struct {
	MinimalHopCount int
	HopCount        int
	Hostkey         []byte

	ListenAddr string
	APIAddr    string
	AuthAddr   string

	Peers    map[p2pnet.Identity]string
	Sessions map[uint32]p2pnet.Identity
	Tunnels  map[uint32]*Tunnel
}

func New(conf *cfg.Configurations) (*Onion, error) {

	var mod *Onion
	var priv *rsa.PrivateKey
	var hostkey []byte
	var err error
	var hostkeyPath string

	mod = &Onion{}
	conf.Init(&mod.MinimalHopCount, ModuleToken, MinHopToken, DefaultMinHop)
	conf.Init(&mod.HopCount, ModuleToken, HopCountToken, DefaultHopCount)
	conf.Init(&mod.ListenAddr, ModuleToken, ListenAddrToken, DefaultListenAddr)
	conf.Init(&mod.APIAddr, ModuleToken, ApiAddrToken, DefaultApiAddr)
	conf.Init(&mod.AuthAddr, auth.ModuleToken, auth.ApiAddrToken, auth.DefaultApiAddr)
	conf.Init(&hostkeyPath, "", HostkeyToken, DefaultHostkey)

	if priv, err = auth.ReadPEMPrivateKey(hostkeyPath); err != nil {
		return nil, err
	}

	if hostkey, err = auth.GetPublicKeyAsDER(&priv.PublicKey); err != nil {
		return nil, err
	}

	mod.Hostkey = hostkey
	mod.Peers = make(map[p2pnet.Identity]string)
	mod.Sessions = make(map[uint32]p2pnet.Identity)
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
		return o.handleTunnelBuild(source, &m)
	case msg.AuthSessionIncomingHS1:
		m := message.(msg.AuthSessionIncomingHS1)
		return o.handleIncomingHS1(source, &m)
	case msg.AuthHandshake2:
		m := message.(msg.AuthHandshake2)
		return o.handleHandshake2(source, &m)
	default:
		return p2pnet.ErrModuleDoesNotHandle
	}
	return nil
}

func (o *Onion) handleTunnelBuild(source net.Conn, m *msg.OnionTunnelBuild) error {

	var port int
	var host net.IP
	var hostport string
	var hostkey []byte
	var tunnelReady *msg.OnionTunnelReady
	var err error

	port = int(m.Port)
	host = net.IP(m.IPAddr)
	hostport = net.JoinHostPort(host.String(), strconv.Itoa(port))
	hostkey = m.DstHostkey

	if tunnelReady, err = o.BuildTunnel(0, hostport, hostkey); err != nil {
		return err
	}
	return msg.Send(source, tunnelReady)
}

func (o *Onion) handleIncomingHS1(source net.Conn, m *msg.AuthSessionIncomingHS1) error {

	var response *msg.AuthSessionHS2
	var sessionId uint32
	var payload []byte
	var err error

	if response, err = o.requestHandshake2(m); err != nil {
		return err
	}

	// Save the session Id and source
	sessionId = response.SessionId
	o.storeSession(sessionId, m.Hostkey)
	o.storeIdentity(source.RemoteAddr().String(), m.Hostkey)

	// Extract the payload.
	payload = response.HandshakePayload

	// Forward the payload to the source.
	if _, err = source.Write(payload); err != nil {
		return err
	}

	return nil
}

func (o *Onion) handleHandshake2(source net.Conn, m msg.Message) error {

	return forwardTo(o.AuthAddr, m)
}

func (o *Onion) handleUnknown(source net.Conn, m msg.Message) error {
	fmt.Printf("Unhandled message type:%v\n", msg.Identifier(m.TypeId()))
	return p2pnet.ErrModuleDoesNotHandle
}
