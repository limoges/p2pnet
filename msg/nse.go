package msg

import (
	"encoding/binary"
	"fmt"
)

const (
	NSE_QUERY    = 520
	NSE_ESTIMATE = 521
	// Reserved up to 539.
)

type NSEQuery struct {
	// This is empty.
}

func (m NSEQuery) String() string {
	return fmt.Sprintf("NSEQuery{}")
}

func (m *NSEQuery) MinimumLength() int {
	return 0
}

func (m *NSEQuery) PayloadLength() int {
	return m.MinimumLength()
}

func (m *NSEQuery) UnmarshalBinary(data []byte) error {
	return nil
}

func (m *NSEQuery) MarshalBinary() (data []byte, err error) {

	payloadBuf := make([]byte, 0)
	return createMessage(NSE_QUERY, payloadBuf), nil
}

type NSEEstimate struct {
	EstimatePeers        uint32
	EstimateStdDeviation uint32
}

func (m NSEEstimate) String() string {
	return fmt.Sprintf(
		"NSEEstimate{EstimatePeers:%v, EstimateStdDeviation:%v}",
		m.EstimatePeers, m.EstimateStdDeviation)
}

func (m *NSEEstimate) MinimumLength() int {
	return 8
}

func (m *NSEEstimate) PayloadLength() int {
	return m.MinimumLength()
}

func (m *NSEEstimate) UnmarshalBinary(data []byte) error {

	if len(data) < m.MinimumLength() {
		return ErrDataTooShort
	}

	m.EstimatePeers = binary.BigEndian.Uint32(data[:4])
	m.EstimateStdDeviation = binary.BigEndian.Uint32(data[4:8])
	return nil
}

func (m *NSEEstimate) MarshalBinary() (data []byte, err error) {

	payloadBuf := make([]byte, m.PayloadLength())

	binary.BigEndian.PutUint32(payloadBuf[:4], m.EstimatePeers)
	binary.BigEndian.PutUint32(payloadBuf[4:8], m.EstimateStdDeviation)
	return createMessage(NSE_ESTIMATE, payloadBuf), nil
}
