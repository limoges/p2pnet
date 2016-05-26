package msg

import (
	"encoding/binary"
	"fmt"
)

// Header represents the protocol's common message header format.
// It is defined to contain message's size in the first 2 bytes,
// followed by the message's type in the following 2 bytes.
type Header struct {
	Size uint16
	Type uint16
}

func (h Header) String() string {
	return fmt.Sprintf("Header{Size:%v, Type:%v}", h.Size, h.Type)
}

func (h Header) MinimumLength() int {
	return HeaderLength
}

func (h Header) PayloadSize() int {
	return int(h.Size) - HeaderLength
}

func (h *Header) UnmarshalBinary(data []byte) error {

	if len(data) < h.MinimumLength() {
		return ErrDataTooShort
	}

	h.Size = binary.BigEndian.Uint16(data[:2])
	h.Type = binary.BigEndian.Uint16(data[2:4])
	return nil
}

func (h Header) MarshalBinary() (data []byte, err error) {

	data = make([]byte, HeaderLength)

	binary.BigEndian.PutUint16(data[0:2], h.Size)
	binary.BigEndian.PutUint16(data[2:4], h.Type)

	return data, nil
}
