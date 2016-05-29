package onion

import (
	"fmt"
	"net"

	"github.com/limoges/p2pnet"
	"github.com/limoges/p2pnet/cfg"
	"github.com/limoges/p2pnet/msg"
)

const (
	moduleToken       = "ONION_FORWARDING"
	minHopToken       = "min_hop_count"
	defaultMinHop     = 2
	hopCountToken     = "hop_count"
	defaultHopCount   = 5
	listenAddrToken   = "listen_address"
	defaultListenAddr = "127.0.0.1:7014"
	apiAddrToken      = "api_address"
	defaultApiAddr    = "127.0.0.1:7004"
)

type Onion struct {
	MinimalHopCount int
	HopCount        int

	ListenAddr string
	APIAddr    string

	Peers    map[string]string
	Sessions map[p2pnet.SessionId]string
}

func New(conf *cfg.Configurations) (onion *Onion, err error) {

	onion = &Onion{}
	conf.Init(&onion.MinimalHopCount, moduleToken, minHopToken, defaultMinHop)
	conf.Init(&onion.HopCount, moduleToken, hopCountToken, defaultHopCount)
	conf.Init(&onion.ListenAddr, moduleToken, listenAddrToken, defaultListenAddr)
	conf.Init(&onion.APIAddr, moduleToken, apiAddrToken, defaultApiAddr)

	onion.Peers = make(map[string]string)
	onion.Sessions = make(map[p2pnet.SessionId]string)
	return onion, nil
}

func (o *Onion) Name() string {
	return moduleToken
}

func (o *Onion) Addresses() (APIAddr, P2PAddr string) {
	return o.APIAddr, o.ListenAddr
}

func (o *Onion) Run() error {

	select {}
}

func (o *Onion) Handle(source net.Conn, message msg.Message) error {
	switch message.(type) {
	case msg.OnionTunnelBuild, *msg.OnionTunnelBuild:
		m := message.(msg.OnionTunnelBuild)
		fmt.Println(m)
	default:
		fmt.Printf("Unhandled message:%v\n", message)
		return p2pnet.ErrModuleDoesNotHandle
	}
	return nil
}

func (o *Onion) BuildTunnel(destAddr net.IP, port int, destHostkey []byte) error {
	return nil
}

func (o *Onion) bootstrapNewPeers() {
}
