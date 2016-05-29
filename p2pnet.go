package p2pnet

import (
	"fmt"
	"net"
)

type Peer struct {
	IPAddr  []byte
	HostKey []byte
}

func BuildConnexionIdentityToken(conn net.Conn) string {
	return fmt.Sprintf("%v/%v", conn.LocalAddr(), conn.RemoteAddr())
}
