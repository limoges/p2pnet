package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"

	"github.com/limoges/p2pnet/auth"
	"github.com/limoges/p2pnet/msg"
)

var (
	sourceAddress      string
	destinationAddress string
)

func init() {
	const (
		defaultAddr = "127.0.0.1:7031"
		usageAddr   = "the address to try to connect to"
	)
	flag.StringVar(&sourceAddress, "s", defaultAddr, usageAddr)
	flag.StringVar(&destinationAddress, "d", defaultAddr, usageAddr)
}

func main() {

	var conn net.Conn
	var err error
	var portStr, hostStr string
	var ips []net.IP
	var ip net.IP
	var hostkey []byte
	var port int

	flag.Parse()

	fmt.Printf("Attempting to build a tunnel between %v and %v\n",
		sourceAddress, destinationAddress)

	// First, split the address into the IP/Host and Port parts.
	if hostStr, portStr, err = net.SplitHostPort(destinationAddress); err != nil {
		fmt.Println(err)
		return
	}

	// Then, lookup the IP in case we have a hostname.
	if ips, err = net.LookupIP(hostStr); err != nil || len(ips) == 0 {
		fmt.Println(err)
		return
	}
	ip = ips[0]

	// Then, read the hostkey of the supposed peer.
	if hostkey, err = ReadSampleHostkey("./samples/sample1.pem"); err != nil {
		fmt.Println(err)
		return
	}

	if port, err = strconv.Atoi(portStr); err != nil {
		fmt.Println(err)
		return
	}

	// Build the OnionTunnelBuild packet.
	tunnelBuild := msg.OnionTunnelBuild{
		Port:       uint16(port),
		IPAddr:     ip.To16(),
		DstHostkey: hostkey,
	}

	// Try to connect to the live source peer
	if conn, err = net.Dial("tcp", sourceAddress); err != nil {
		fmt.Println(err)
		return
	}

	if err = msg.WriteMessage(conn, tunnelBuild); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Success?")
}

func ReadSampleHostkey(path string) ([]byte, error) {

	var keys *auth.Encryption
	var err error

	if keys, err = auth.ReadKeys(path); err != nil {
		return nil, err
	}

	return keys.Hostkey, nil
}
