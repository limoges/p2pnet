package msg

import (
	"bytes"
	"encoding/binary"
)

const (
	NSE_QUERY    = 520
	NSE_ESTIMATE = 521
	// Reserved up to 539.
)

type NSEQuery struct {
	// This is empty.
}

func (m NSEQuery) TypeId() uint16 {
	return NSE_QUERY
}

func NewNSEQuery(data []byte) (NSEQuery, error) {

	return NSEQuery{}, nil
}

type NSEEstimate struct {
	Peers     uint32
	Deviation uint32
}

func (m NSEEstimate) TypeId() uint16 {
	return NSE_ESTIMATE
}

func NewNSEEstimate(data []byte) (NSEEstimate, error) {

	var m NSEEstimate
	var buf *bytes.Reader
	var err error

	m = NSEEstimate{}
	buf = bytes.NewReader(data)
	if err = binary.Read(buf, binary.BigEndian, &m.Peers); err != nil {
		return m, err
	}

	if err = binary.Read(buf, binary.BigEndian, &m.Deviation); err != nil {
		return m, err
	}

	return m, nil
}
