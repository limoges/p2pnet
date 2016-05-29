package p2pnet

import (
	"net"
)

type Module interface {
	Run() error
	ListenAPI(ln net.Listener)
}
