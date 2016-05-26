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
	CacheSize         int
	MaximumConnexions int
	BootstrapAddress  string
	ListenAddress     string
	ApiAddress        string
}

func New(conf cfg.Configurations) (module Gossip, err error) {

	// TODO: Add defaults to these configurations.
	conf.Init(&module.CacheSize, sectionToken,
		"cache_size", 50)
	conf.Init(&module.MaximumConnexions, sectionToken,
		"max_connections", 3)
	conf.Init(&module.BootstrapAddress, sectionToken,
		"bootstrapper", "fulcrum.net.in.tum.de:6001")
	conf.Init(&module.ListenAddress, sectionToken,
		"listen_address", "127.0.0.1:6001")
	conf.Init(&module.ApiAddress, sectionToken,
		"api_address", "127.0.0.1:7001")

	return module, nil
}
