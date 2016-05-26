package msg

import (
	"encoding/binary"
	"fmt"
)

type RPSQuery struct {
	// This is empty.
}

func (m RPSQuery) String() string {
	return fmt.Sprintf("RPSQuery{}")
}

func (m *RPSQuery) MinimumLength() int {
	return 0
}

func (m *RPSQuery) PayloadLength() int {
	return m.MinimumLength()
}

func (m *RPSQuery) UnmarshalBinary(data []byte) error {
	return nil
}

func (m *RPSQuery) MarshalBinary() (data []byte, err error) {

	payloadBuf := make([]byte, 0)
	return createMessage(RPS_QUERY, payloadBuf), nil
}

type RPSPeer struct {
	Port             uint16
	reserved         uint16
	IPAddress        []byte
	PeerHostKeyInDER []byte
}

func (m RPSPeer) String() string {
	return fmt.Sprintf("RPSPeer{Port:%v, IPAddress:%v, PeerHostKeyInDER:%v}", m.Port,
		m.IPAddress, m.PeerHostKeyInDER)
}

func (m *RPSPeer) MinimumLength() int {
	return 4
}

func (m *RPSPeer) PayloadLength() int {

	return m.MinimumLength() + IPLength
}

func (m *RPSPeer) UnmarshalBinary(data []byte) error {

	if len(data) < m.MinimumLength() {
		return ErrDataTooShort
	}

	m.Port = binary.BigEndian.Uint16(data[:2])
	m.IPAddress = make([]byte, IPLength)
	copy(m.IPAddress, data[4:4+IPLength])
	m.PeerHostKeyInDER = make([]byte, len(data[4+IPLength:]))
	copy(m.PeerHostKeyInDER, data[4+IPLength:])
	return nil
}

func (m RPSPeer) MarshalBinary() (data []byte, err error) {

	payloadBuf := make([]byte, m.PayloadLength())
	binary.BigEndian.PutUint16(payloadBuf[:2], m.Port)

	copy(payloadBuf[4:4+IPLength], m.IPAddress)
	copy(payloadBuf[4+IPLength:], m.PeerHostKeyInDER)
	return createMessage(RPS_PEER, payloadBuf), nil
}
