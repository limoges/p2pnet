package gossip

import (
	"net"

	"github.com/limoges/p2pnet"
	"github.com/limoges/p2pnet/cfg"
	"github.com/limoges/p2pnet/msg"
)

type Gossip struct {
	CacheSize     int
	MaxConns      int
	BootstrapAddr string
	ListenAddr    string
	APIAddr       string
}

const (
	PeerInformation   = 100
	cacheSizeToken    = "cache_size"
	DefaultCacheSize  = 50
	maxConnsToken     = "max_connections"
	DefaultMaxConns   = 3
	bootstrapToken    = "bootstrapper"
	DefaultBootstrap  = "fulcrum.net.in.tum.de:6001"
	listenAddrToken   = "listen_address"
	DefaultListenAddr = "127.0.0.1:7001"
	apiAddrToken      = "api_address"
	DefaultApiAddr    = "127.0.0.1:6001"
	ModuleToken       = "GOSSIP"
)

func New(conf *cfg.Configurations) (module *Gossip, err error) {

	module = &Gossip{}
	conf.Init(&module.CacheSize, ModuleToken, cacheSizeToken, DefaultCacheSize)
	conf.Init(&module.MaxConns, ModuleToken, maxConnsToken, DefaultMaxConns)
	conf.Init(&module.BootstrapAddr, ModuleToken, bootstrapToken, DefaultBootstrap)
	conf.Init(&module.ListenAddr, ModuleToken, listenAddrToken, DefaultListenAddr)
	conf.Init(&module.APIAddr, ModuleToken, apiAddrToken, DefaultApiAddr)
	return module, nil
}

func (m *Gossip) Name() string {
	return ModuleToken
}

func (m *Gossip) Addresses() (APIAddr, P2PAddr string) {
	return m.APIAddr, m.ListenAddr
}

func (m *Gossip) Run() error {

	select {
	// Let's the process idle.
	}
}

func (m *Gossip) Handle(source net.Conn, message msg.Message) error {

	switch message.(type) {
	default:
		return p2pnet.ErrModuleDoesNotHandle
	}
}
