package main

import (
	"crypto/rsa"
	"errors"
	"flag"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/limoges/p2pnet/auth"
	"github.com/limoges/p2pnet/cfg"
	"github.com/limoges/p2pnet/client"
	"github.com/limoges/p2pnet/msg"
)

const (
	ConfigDir = "configs/"
)

const (
	numberOfPeers = 5
)

type Peer struct {
	name   string
	config *cfg.Configurations

	Hostkey []byte
	Client  *client.Client
}

func NewPeer(i int) *Peer {

	var name string
	name = fmt.Sprintf("peer%v", i)
	peer := &Peer{name: name}
	return peer
}

func (p *Peer) filename() string {
	return ConfigDir + p.name + ".ini"
}

func (p *Peer) loadClient() error {

	var filename string
	var c *client.Client
	var err error

	filename = p.filename()

	if c, err = client.New(filename); err != nil {
		fmt.Println(err)
		return err
	} else {
		p.Client = c
		go p.Client.Run()
	}
	return nil
}

func splitparse(hostport string) (ip net.IP, port int, err error) {

	var hostString, portString string
	var ips []net.IP

	if hostString, portString, err = net.SplitHostPort(hostport); err != nil {
		return nil, 0, err
	}

	if ips, err = net.LookupIP(hostString); err != nil {
		return nil, 0, err
	}
	ip = ips[0]

	if port, err = strconv.Atoi(portString); err != nil {
		return nil, 0, err
	}

	return ip, port, nil
}

func StartSessionBetweenPeers(source, target *Peer) error {

	var conn net.Conn
	var message *msg.OnionTunnelBuild
	var err error
	var pub *rsa.PublicKey
	var targetHostkey []byte
	var sourceAddr, targetAddr string
	var ip net.IP
	var port int
	var response msg.Message

	pub = &target.Client.ModAuth.PrivateKey.PublicKey

	sourceAddr = source.Client.ModOnion.ListenAddr
	targetAddr = target.Client.ModOnion.ListenAddr

	fmt.Printf("Source address %v\n", sourceAddr)
	fmt.Printf("Target address %v\n", targetAddr)

	if ip, port, err = splitparse(targetAddr); err != nil {
		return err
	}

	if targetHostkey, err = auth.MarshalPublicKey(pub); err != nil {
		return err
	}

	fmt.Printf("Hostkey length in bytes is %v\n", len(targetHostkey))

	if conn, err = net.Dial("tcp", sourceAddr); err != nil {
		return err
	}

	message = &msg.OnionTunnelBuild{
		Port:       uint16(port),
		IPAddr:     ip.To16(),
		DstHostkey: targetHostkey,
	}

	if response, err = msg.SendReceive(conn, message); err != nil {
		return err
	}

	if response.TypeId() != msg.ONION_TUNNEL_READY {
		return errors.New("Could not complete tunnel.")
	}

	fmt.Println("Sending data through...")
	return nil
}

func main() {

	const numberOfPeers = 2
	var err error

	flag.Parse()

	peers := make([]*Peer, 0, numberOfPeers)

	for i := 0; i < numberOfPeers; i++ {
		peer := NewPeer(i)
		peer.loadClient()
		peers = append(peers, peer)
	}

	time.Sleep(3 * time.Second)

	if err = StartSessionBetweenPeers(peers[0], peers[1]); err != nil {
		fmt.Println(err)
	}

	select {}
}
