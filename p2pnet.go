package p2pnet

import (
	"fmt"
	"net"
)

type SessionId uint32

type Session struct {
	ID      SessionId
	Partner Peer
}

func BuildConnexionIdentityToken(conn net.Conn) string {
	return fmt.Sprintf("%v/%v", conn.LocalAddr(), conn.RemoteAddr())
}
