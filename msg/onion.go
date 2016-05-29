package msg

import (
	"encoding/binary"
	"fmt"
)

const (
	ONION_TUNNEL_BUILD    = 560
	ONION_TUNNEL_READY    = 561
	ONION_TUNNEL_INCOMING = 562
	ONION_TUNNEL_DESTROY  = 563
	ONION_TUNNEL_DATA     = 564
	ONION_ERROR           = 565
	ONION_COVER           = 566
	// Reserved up to 599.
)

type OnionTunnelBuild struct {
	reserved                uint16
	Port                    uint16
	IPAddr                  []byte
	DestinationHostKeyInDER []byte
}

func (m OnionTunnelBuild) String() string {
	return fmt.Sprintf(
		"OnionTunnelBuild{Port:%v, IPAddr:%v, DestinationHostKeyInDER:%v",
		m.Port, m.IPAddr, m.DestinationHostKeyInDER)
}

func (m OnionTunnelBuild) MinimumLength() int {
	return 4
}

func (m OnionTunnelBuild) PayloadLength() int {
	return m.MinimumLength() + len(m.IPAddr) + len(m.DestinationHostKeyInDER)
}

func (m *OnionTunnelBuild) UnmarshalBinary(data []byte) error {

	if len(data) < m.MinimumLength() {
		return ErrDataTooShort
	}

	m.Port = binary.BigEndian.Uint16(data[2:4])
	m.IPAddr = make([]byte, IPLength)
	copy(m.IPAddr, data[4:4+IPLength])
	m.DestinationHostKeyInDER = make([]byte, len(data[4+IPLength:]))
	copy(m.DestinationHostKeyInDER, data[4+IPLength:])
	return nil
}

func (m OnionTunnelBuild) MarshalBinary() (data []byte, err error) {

	payloadBuf := make([]byte, m.PayloadLength())
	binary.BigEndian.PutUint16(payloadBuf[2:4], m.Port)

	copy(payloadBuf[4:4+IPLength], m.IPAddr)
	copy(payloadBuf[4+IPLength:], m.DestinationHostKeyInDER)
	return createMessage(ONION_TUNNEL_BUILD, payloadBuf), nil
}

type OnionTunnelReady struct {
	TunnelID                uint32
	DestinationHostKeyInDER []byte
}

func (m OnionTunnelReady) String() string {
	return fmt.Sprintf(
		"OnionTunnelReady{TunnelID:%v, DestinationHostKeyInDER:%v}",
		m.TunnelID, m.DestinationHostKeyInDER)
}

func (m OnionTunnelReady) MinimumLength() int {
	return 4
}

func (m OnionTunnelReady) PayloadLength() int {
	return m.MinimumLength() + len(m.DestinationHostKeyInDER)
}

func (m *OnionTunnelReady) UnmarshalBinary(data []byte) error {

	if len(data) < m.MinimumLength() {
		return ErrDataTooShort
	}

	m.TunnelID = binary.BigEndian.Uint32(data[:4])
	m.DestinationHostKeyInDER = make([]byte, len(data[4:]))
	copy(m.DestinationHostKeyInDER, data[4:])
	return nil
}

func (m OnionTunnelReady) MarshalBinary() (data []byte, err error) {

	payloadBuf := make([]byte, m.PayloadLength())

	binary.BigEndian.PutUint32(data[:4], m.TunnelID)
	copy(payloadBuf[4:], m.DestinationHostKeyInDER)
	return createMessage(ONION_TUNNEL_READY, payloadBuf), nil
}

type OnionTunnelIncoming struct {
	TunnelID           uint32
	SourceHostKeyInDER []byte
}

func (m OnionTunnelIncoming) String() string {
	return fmt.Sprintf(
		"OnionTunnelIncoming{TunnelID:%v, SourceHostKeyInDER:%v}",
		m.TunnelID, m.SourceHostKeyInDER)
}

func (m OnionTunnelIncoming) MinimumLength() int {
	return 4
}

func (m OnionTunnelIncoming) PayloadLength() int {
	return m.MinimumLength() + len(m.SourceHostKeyInDER)
}

func (m *OnionTunnelIncoming) UnmarshalBinary(data []byte) error {

	if len(data) < m.MinimumLength() {
		return ErrDataTooShort
	}

	m.TunnelID = binary.BigEndian.Uint32(data[:4])
	m.SourceHostKeyInDER = make([]byte, len(data[4:]))
	copy(m.SourceHostKeyInDER, data[4:])
	return nil
}

func (m OnionTunnelIncoming) MarshalBinary() (data []byte, err error) {

	payloadBuf := make([]byte, m.PayloadLength())

	binary.BigEndian.PutUint32(data[:4], m.TunnelID)
	copy(payloadBuf[4:], m.SourceHostKeyInDER)
	return createMessage(ONION_TUNNEL_INCOMING, payloadBuf), nil
}

type OnionTunnelDestroy struct {
	TunnelID uint32
}

func (m OnionTunnelDestroy) String() string {
	return fmt.Sprintf(
		"OnionTunnelDestroy{TunnelID:%v}", m.TunnelID)
}

func (m OnionTunnelDestroy) MinimumLength() int {
	return 4
}

func (m OnionTunnelDestroy) PayloadLength() int {
	return m.MinimumLength()
}

func (m *OnionTunnelDestroy) UnmarshalBinary(data []byte) error {

	if len(data) < m.MinimumLength() {
		return ErrDataTooShort
	}

	m.TunnelID = binary.BigEndian.Uint32(data[:4])
	return nil
}

func (m OnionTunnelDestroy) MarshalBinary() (data []byte, err error) {

	payloadBuf := make([]byte, m.PayloadLength())

	binary.BigEndian.PutUint32(data[:4], m.TunnelID)
	return createMessage(ONION_TUNNEL_DESTROY, payloadBuf), nil
}

type OnionTunnelData struct {
	TunnelID uint32
	Data     []byte
}

func (m OnionTunnelData) String() string {
	return fmt.Sprintf(
		"OnionTunnelData{TunnelID:%v, Data:%v}",
		m.TunnelID, m.Data)
}

func (m OnionTunnelData) MinimumLength() int {
	return 4
}

func (m OnionTunnelData) PayloadLength() int {
	return m.MinimumLength() + len(m.Data)
}

func (m *OnionTunnelData) UnmarshalBinary(data []byte) error {

	if len(data) < m.MinimumLength() {
		return ErrDataTooShort
	}

	m.TunnelID = binary.BigEndian.Uint32(data[:4])
	m.Data = make([]byte, len(data[4:]))
	copy(m.Data, data[4:])
	return nil
}

func (m OnionTunnelData) MarshalBinary() (data []byte, err error) {

	payloadBuf := make([]byte, m.PayloadLength())

	binary.BigEndian.PutUint32(data[:4], m.TunnelID)
	copy(payloadBuf[4:], m.Data)
	return createMessage(ONION_TUNNEL_DATA, payloadBuf), nil
}

type OnionError struct {
	RequestType uint16
	reserved    uint16
	TunnelID    uint32
}

func (m OnionError) String() string {
	return fmt.Sprintf(
		"OnionError{RequestType:%v, TunnelID:%v}",
		m.RequestType, m.TunnelID)
}

func (m OnionError) MinimumLength() int {
	return 8
}

func (m OnionError) PayloadLength() int {
	return m.MinimumLength()
}

func (m *OnionError) UnmarshalBinary(data []byte) error {

	if len(data) < m.MinimumLength() {
		return ErrDataTooShort
	}

	m.RequestType = binary.BigEndian.Uint16(data[:2])
	m.TunnelID = binary.BigEndian.Uint32(data[4:8])
	return nil
}

func (m OnionError) MarshalBinary() ([]byte, error) {

	payloadBuf := make([]byte, m.PayloadLength())

	binary.BigEndian.PutUint16(payloadBuf[:2], m.RequestType)
	binary.BigEndian.PutUint32(payloadBuf[4:8], m.TunnelID)
	return createMessage(ONION_ERROR, payloadBuf), nil
}

type OnionCover struct {
	CoverSize uint16
	reserved  uint16
}

func (m OnionCover) String() string {
	return fmt.Sprintf("OnionCover{CoverSize:%v}", m.CoverSize)
}

func (m OnionCover) MinimumLength() int {
	return 4
}

func (m OnionCover) PayloadLength() int {
	return m.MinimumLength()
}

func (m *OnionCover) UnmarshalBinary(data []byte) error {

	if len(data) < m.MinimumLength() {
		return ErrDataTooShort
	}

	m.CoverSize = binary.BigEndian.Uint16(data[:2])
	return nil
}

func (m OnionCover) MarshalBinary() ([]byte, error) {

	payloadBuf := make([]byte, m.PayloadLength())

	binary.BigEndian.PutUint16(payloadBuf[:2], m.CoverSize)
	return createMessage(ONION_COVER, payloadBuf), nil
}
