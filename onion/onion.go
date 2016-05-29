package onion

import (
	"github.com/limoges/p2pnet/cfg"
)

type Onion struct {
	MinimalHopCount int
	HopCount        int
	IPAddr          []byte
}

func New(conf *cfg.Configurations) (onion *Onion, err error) {

	onion = &Onion{}
	return onion
}
