package main

import (
	"flag"
	"fmt"
	"github.com/limoges/p2pnet/msg"
	"net"
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
	fmt.Println("Successfully connected.")

	announce := &msg.GossipAnnounce{
		TTL:      1,
		DataType: 17,
		Data:     []byte{1, 1, 1},
	}

	notify := &msg.GossipNotify{
		DataType: 1,
	}

	notification := &msg.GossipNotification{
		HeaderID: 32,
		DataType: 3,
		Data:     []byte{123, 123, 123},
	}

	validation := &msg.GossipValidation{
		MessageID: 42,
		Valid:     true,
	}

	nseQuery := &msg.NSEQuery{}

	nseEstimate := &msg.NSEEstimate{
		EstimatePeers:        1,
		EstimateStdDeviation: 1,
	}

	rpsQuery := &msg.RPSQuery{}

	rpsPeer := &msg.RPSPeer{
		Port:             80,
		IPAddress:        net.IPv4(192, 168, 0, 1),
		PeerHostKeyInDER: []byte{128, 128, 128},
	}

	messages := []msg.Message{
		announce,
		notify,
		notification,
		validation,
		nseQuery,
		nseEstimate,
		rpsQuery,
		rpsPeer,
	}

	for _, m := range messages {
		fmt.Println(m)
		buf, _ := m.MarshalBinary()
		fmt.Println(buf)
		_, err := conn.Write(buf)
		if err != nil {
			fmt.Println(err)
		}
	}
}
