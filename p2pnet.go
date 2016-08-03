package p2pnet

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
)

// Represents the identity of a peer. It corresponds to the SHA256 checksum of
// the peer's hostkey.
type Identity string
type Hostkey []byte

func GetIdentity(hostkey []byte) Identity {

	var sum [sha256.Size]byte
	sum = sha256.Sum256(hostkey)
	return Identity(hex.EncodeToString(sum[:]))
}

type SessionId uint32

type Session struct {
	ID      SessionId
	Partner Peer
}

func BuildConnexionIdentityToken(conn net.Conn) string {
	return fmt.Sprintf("%v/%v", conn.LocalAddr(), conn.RemoteAddr())
}

type Peer struct {
	Port    uint16
	IPAddr  []byte
	Hostkey []byte
}

func (p Peer) MarshalJSON() ([]byte, error) {
	return json.Marshal(p)
}

func (p *Peer) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, p)
}
