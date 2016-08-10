package client

import (
	"github.com/limoges/p2pnet"
	"github.com/limoges/p2pnet/auth"
	"github.com/limoges/p2pnet/cfg"
	"github.com/limoges/p2pnet/onion"
)

type Client struct {
	ListenAddr string
	Config     *cfg.Configurations

	ModAuth  *auth.Auth
	ModOnion *onion.Onion
	// ModRPS    *rps.RPS
	// ModNSE    *nse.NSE
	// ModGossip *gossip.Gossip
	Modules []p2pnet.Module
}

func New(filename string) (*Client, error) {

	var config *cfg.Configurations
	var client *Client
	var err error

	client = &Client{}

	if config, err = cfg.New(filename); err != nil {
		return nil, err
	}
	client.Config = config

	if client.ModAuth, err = auth.New(config); err != nil {
		return nil, err
	}

	if client.ModOnion, err = onion.New(config); err != nil {
		return nil, err
	}

	//if client.modRPS, err = rps.New(config); err != nil {
	//	return nil, err
	//}

	//if client.modNSE, err = nse.New(config); err != nil {
	//	return nil, err
	//}

	//if client.modGossip, err = gossip.New(config); err != nil {
	//	return nil, err
	//}

	client.Modules = []p2pnet.Module{
		client.ModAuth,
		client.ModOnion,
		// client.modRPS,
		// client.modNSE,
		// client.modGossip,
	}

	return client, nil
}

func (c *Client) Run() error {

	// peersListener, err := net.Listen("tcp", c.ListenAddr)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return err
	// }

	for _, module := range c.Modules {
		go p2pnet.Run(module)
	}

	select {}
}
