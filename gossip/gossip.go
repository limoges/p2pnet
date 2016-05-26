package gossip

import (
	"fmt"
	"github.com/limoges/p2pnet/cfg"
	"github.com/limoges/p2pnet/msg"
	"io"
	"net"
)

type Gossip struct {
	CacheSize         int
	MaximumConnexions int
	BootstrapAddress  string
	ListenAddress     string
	ApiAddress        string
}

func (g *Gossip) LoadConfiguration(c cfg.Configurations) {

	const (
		sectionToken           = "Gossip"
		cacheSizeToken         = "cache_size"
		degreeToken            = "degree"
		maximumConnexionsToken = "max_connections"
		bootstrapToken         = "bootstrapper"
		listenToken            = "listen_address"
		apiToken               = "api_address"
	)

	if c.File == nil {
		fmt.Printf("Cannot load Gossip configurations. No configuration provided.")
		return
	}

	section, err := c.File.GetSection(sectionToken)

	if err != nil {
		fmt.Println(err)
		return
	}

	g.CacheSize, _ = section.Key(cacheSizeToken).Int()
	g.MaximumConnexions, _ = section.Key(maximumConnexionsToken).Int()
	g.BootstrapAddress = section.Key(bootstrapToken).String()
	g.ListenAddress = section.Key(listenToken).String()
	g.ApiAddress = section.Key(apiToken).String()

	fmt.Println("Loaded Gossip configurations:", g)
}

func (g *Gossip) Run() {

	fmt.Printf("Launching Gossip module on %s\n", g.ApiAddress)

	ln, err := net.Listen("tcp", g.ApiAddress)

	if err != nil {
		fmt.Println(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
		} else {
			go g.Handle(conn)
		}
	}
}

func buildConnexionIdentityToken(conn net.Conn) string {
	return fmt.Sprintf("%v/%v", conn.LocalAddr(), conn.RemoteAddr())
}

func (g *Gossip) Handle(conn net.Conn) {

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
			fmt.Println("WTF")
		}

		packetBuf := make([]byte, 0)
		packetBuf = append(packetBuf, headerBuf...)
		packetBuf = append(packetBuf, payloadBuf...)
		fmt.Println(packetBuf)

		// We now have the header and payload, just build the message.
		msg.NewMessage(header, payloadBuf)
	}
}
