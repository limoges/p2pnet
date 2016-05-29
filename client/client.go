package client

import (
	"fmt"
	"github.com/limoges/p2pnet/auth"
	"github.com/limoges/p2pnet/cfg"
	"github.com/limoges/p2pnet/msg"
	"io"
	"net"
)

type Client struct {
	APIAddr    string
	ListenAddr string
	Auth       *auth.Auth
}

func New(conf *cfg.Configurations) (client *Client, err error) {

	client = &Client{}
	conf.Init(&client.APIAddr, "GOSSIP", "api_address", "127.0.0.1:7011")
	conf.Init(&client.ListenAddr, "GOSSIP", "listen_address", "127.0.0.1:6011")

	// Load hostkey into memory
	auth, err := auth.New(conf)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	client.Auth = auth

	return client, nil
}

func (c *Client) Run() error {

	apiListener, err := net.Listen("tcp", c.APIAddr)
	if err != nil {
		fmt.Println(err)
		return err
	}

	peersListener, err := net.Listen("tcp", c.ListenAddr)
	if err != nil {
		fmt.Println(err)
		return err
	}

	go c.listenAPI(apiListener)
	go c.listenPeers(peersListener)

	for {

		// Start the active stuff, like going checking on the bootstrapper.
	}
}

func (c *Client) listenAPI(ln net.Listener) {

	fmt.Printf("Launched API listener on %v\n", ln.Addr())
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
		} else {
			go c.handle(conn)
		}
	}
}

func (c *Client) listenPeers(ln net.Listener) {

	fmt.Printf("Launched peer listener on %v\n", ln.Addr())
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
		} else {
			go c.handle(conn)
		}
	}
}
