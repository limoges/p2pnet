package p2pnet

import "encoding/json"

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
