package gossip

import (
	"fmt"
	"github.com/limoges/p2pnet/cfg"
	"net"
)

const (
	sectionToken           = "Gossip"
	cacheSizeToken         = "cache_size"
	degreeToken            = "degree"
	maximumConnexionsToken = "max_connections"
	bootstrapToken         = "bootstrapper"
	listenToken            = "listen_address"
	apiToken               = "api_address"
)

type Gossip struct {
	CacheSize     int
	MaxConns      int
	BootstrapAddr string
	ListenAddr    string
	APIAddr       string
}

type Subscriber struct {
	Addr     []byte
	DataType int
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
)

func New(conf cfg.Configurations) (module Gossip, err error) {

	module := &Gossip{}
	conf.Init(&module.CacheSize, sectionToken, cacheSizeToken, DefaultCacheSize)
	conf.Init(&module.MaxConns, sectionToken, maxConnsToken, DefaultMaxConns)
	conf.Init(&module.BootstrapAddr, sectionToken, bootstrapToken, DefaultBootstrap)
	conf.Init(&module.ListenAddr, sectionToken, listenAddrToken, DefaultListenAddr)
	conf.Init(&module.APIAddr, sectionToken, apiAddrToken, DefaultApiAddr)
	return module, nil
}

func (g *Gossip) Run() {

	apiListener, err := net.Listen("tcp", g.APIAddr)
	if err != nil {
		fmt.Println(err)
	}

	listener, err := net.Listen("tcp", g.ListenAddr)
	if err != nil {
		fmt.Println(err)
	}

	go g.listenAPI(apiListener)
	go g.listen(listener)

	for {
	}
}

func (g *Gossip) listenAPI(ln net.Listener) {

	fmt.Printf("Launched Gossip API Listener on %v\n", ln.Addr())

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
		} else {
			go g.handle(conn)
		}
	}
}

func (g *Gossip) handleAPI(conn net.Conn) {

}

func (g *Gossip) listen(ln net.Listener) {
}
