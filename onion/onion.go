package onion

import (
	"github.com/limoges/p2pnet/cfg"
)

type Onion struct {
	MinimalHopCount int
	HopCount        int
	IPAddr          []byte

	Peers    map[string]string
	Sessions [uint32]string
}

func New(conf *cfg.Configurations) (onion *Onion, err error) {

	onion = &Onion{}
	onion.Peers = make(map[string][]byte)
	return onion
}

func (o *Onion) handle(m Message) {

	switch m.(type) {
	case OnionTunnelBuild:
		om := OnionTunnelBuild(m)
		o.buildTunnel(om.Port, om.IPAddr, om.DestinationHostKeyInDER)
	case OnionTunnelReady:
	case OnionTunnelIncoming:
	case OnionTunnelDestroy:
	case OnionTunnelData:
	case OnionError:
	case OnionCover:
	default:
		fmt.Printf("Unhandled message:%v\n", m)
	}
}

func (o *Onion) BuildTunnel(port uint16, addr []byte, hostkey []byte) {

	// Create or update the corresponding entry in our table.
	o.Peers[string(hostkey)] = fmt.Sprintf("%v:%v", IP(addr), port)

	// Send a request to the auth module to start a session.
	m := AuthSessionStart{
		Hostkey: hostkey,
	}

	// Send to local auth module.
}
