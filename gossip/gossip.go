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
	PeerInformation = 100
)

func New(conf cfg.Configurations) (module Gossip, err error) {

	// TODO: Add defaults to these configurations.
	conf.Init(&module.CacheSize, sectionToken,
		"cache_size", 50)
	conf.Init(&module.MaxConns, sectionToken,
		"max_connections", 3)
	conf.Init(&module.BootstrapAddr, sectionToken,
		"bootstrapper", "fulcrum.net.in.tum.de:6001")
	conf.Init(&module.ListenAddr, sectionToken,
		"listen_address", "127.0.0.1:6001")
	conf.Init(&module.APIAddr, sectionToken,
		"api_address", "127.0.0.1:7001")

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
