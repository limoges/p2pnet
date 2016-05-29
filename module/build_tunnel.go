package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/limoges/p2pnet/msg"
)

var address string

func init() {
	const (
		defaultAddr = "127.0.0.1:7001"
		usageAddr   = "the address to try to connect to"
	)
	flag.StringVar(&address, "a", defaultAddr, usageAddr)
}

func main() {

	flag.Parse()

	fmt.Printf("Testing connexion: %v\n", address)

	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println(err)
		return
	}

	sessionStart := msg.AuthSessionStart{
		Hostkey: []byte{255, 255, 255, 255},
	}

	sessionHS1 := msg.AuthSessionHS1{
		SessionId:        uint32(123456),
		HandshakePayload: []byte{255, 254, 253, 252},
	}

	messages := []msg.Message{
		sessionStart,
		sessionHS1,
	}

	for _, m := range messages {
		err := msg.WriteMessage(conn, m)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
