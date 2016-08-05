package rps

import (
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/limoges/p2pnet"
	"github.com/limoges/p2pnet/cfg"
	"github.com/limoges/p2pnet/msg"
)

const (
	moduleToken          = "RPS"
	listenAddrToken      = "listen_address"
	defaultListenAddr    = "127.0.0.1:7021"
	apiAddrToken         = "api_address"
	defaultApiAddr       = "127.0.0.1:7022"
	bootstrapAddrToken   = "bootstrap_address"
	defaultBootstrapAddr = "127.0.0.1:6001"
)

type RPS struct {
	ListenAddr    string
	APIAddr       string
	BootstrapAddr string

	PeersByIP      map[string]p2pnet.Peer
	PeersByHostkey map[string]p2pnet.Peer
}

func New(conf *cfg.Configurations) (*RPS, error) {

	rps := &RPS{}
	conf.Init(&rps.ListenAddr, moduleToken, listenAddrToken, defaultListenAddr)
	conf.Init(&rps.APIAddr, moduleToken, apiAddrToken, defaultApiAddr)
	conf.Init(&rps.BootstrapAddr, moduleToken, bootstrapAddrToken, defaultBootstrapAddr)

	rps.PeersByIP = make(map[string]p2pnet.Peer)
	rps.PeersByHostkey = make(map[string]p2pnet.Peer)
	return rps, nil
}

func (r *RPS) Name() string {
	return moduleToken
}

func (r *RPS) Addresses() (APIAddr, P2PAddr string) {
	return r.APIAddr, r.ListenAddr
}

func (r *RPS) Run() error {

	select {}
}

func (r *RPS) Handle(source net.Conn, message msg.Message) error {

	switch message.(type) {
	case msg.RPSQuery:
		r.ReplyWithRandomPeer(source)
	default:
		return p2pnet.ErrModuleDoesNotHandle
	}
	return nil
}

func (r *RPS) getRandomPeer() (p2pnet.Peer, error) {

	seed := time.Now().UnixNano()
	source := rand.NewSource(seed)
	rng := rand.New(source)

	peerNo := rng.Intn(len(r.PeersByIP))
	i := 0
	for key := range r.PeersByIP {
		if i == peerNo {
			return r.PeersByIP[key], nil
		}
	}
	return p2pnet.Peer{}, nil
}

func (r *RPS) ReplyWithRandomPeer(source net.Conn) error {

	// Get a random peer from our list of known peers.
	peer, err := r.getRandomPeer()
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Send the random peer information response.
	response := msg.RPSPeer{
		Port:    peer.Port,
		Hostkey: peer.Hostkey,
	}
	copy(response.IPAddr[:], peer.IPAddr)

	return msg.WriteMessage(source, response)
}
