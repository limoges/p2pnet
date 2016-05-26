package client

import (
	"fmt"
	"github.com/limoges/p2pnet/cfg"
	"github.com/limoges/p2pnet/msg"
	"io"
	"net"
)

type Client struct {
	APIAddr    string
	ListenAddr string
}

func buildConnexionIdentityToken(conn net.Conn) string {
	return fmt.Sprintf("%v/%v", conn.LocalAddr(), conn.RemoteAddr())
}

func New(conf *cfg.Configurations) (client *Client, err error) {

	client = &Client{}
	conf.Init(&client.APIAddr, "GOSSIP", "api_address", "127.0.0.1:7011")
	conf.Init(&client.ListenAddr, "GOSSIP", "listen_address", "127.0.0.1:6011")
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

func (c *Client) handle(conn net.Conn) {

	headerBuf := make([]byte, msg.HeaderLength)

	// Build the token we use for logging purposes
	token := buildConnexionIdentityToken(conn)
	fmt.Printf("[%v] New connexion opened by peer\n", token)

	// Run the handling loop, which is supposed to read packets, until
	// the connexion receives the EOF token, signifying that the peer
	// has closed the connexion.
	for {

		// First, we read the message's header.
		bytesRead, err := conn.Read(headerBuf)

		// Check read errors and connexion status.
		if err != nil {

			// EOF error is put up when the connexion has been closed by peer.
			if err == io.EOF {
				fmt.Printf("[%v] Connexion closed by peer.\n", token)
				// We simply stop handling this connexion.
				return
			}

			fmt.Println(err)
			continue
		}

		// Don't bother with the rest if the message is not long enough.
		header := msg.Header{}
		err = header.UnmarshalBinary(headerBuf[:bytesRead])
		if err != nil {
			// We have a non-conforming header
			fmt.Println(err)
			continue
		}

		// Next we read the payload data.
		payloadBuf := make([]byte, header.PayloadSize())

		bytesRead, err = conn.Read(payloadBuf)
		if bytesRead != len(payloadBuf) {
			continue
		}

		packetBuf := make([]byte, 0)
		packetBuf = append(packetBuf, headerBuf...)
		packetBuf = append(packetBuf, payloadBuf...)

		// We now have the header and payload, just build the message.
		message := msg.NewMessage(header, payloadBuf)
		go c.handleMessage(message)
	}
}

func (c *Client) handleMessage(message msg.Message) {
	fmt.Println(message)
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
