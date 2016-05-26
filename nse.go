package main

import ()

type Peer struct {
	Identifier  string
	NetworkAddr string
}

type NSE struct {
	Peers []Peer
}
