package msg

import (
	"bufio"
	"errors"
	"io"
	"log"
)

const (
	HeaderLength = 4
	IPLength     = 16
)

var (
	ErrDataTooShort = errors.New("unmarshal: data is shorter than expected")
)

type Message interface {
	TypeId() uint16
}

// Read a message from the reader.
func ReadMessage(reader io.Reader) (Message, error) {

	var buf *bufio.Reader

	// Create a buffered reader because calls to read are blocking.
	// If we have an ill-formed message, we could be stuck trying to
	// wait for bytes that will never arrive.
	buf = bufio.NewReader(reader)

	if generic, err := ReadGenericMessage(buf); err != nil {
		return nil, err
	} else {
		log.Println("In :", generic)
		m, err := ConvertFromGeneric(generic)
		return m, err
	}
}

// Write a message to the writer.
func WriteMessage(writer io.Writer, m Message) error {

	generic, err := ConvertToGeneric(m)
	if err != nil {
		return err
	}
	log.Println("Out:", generic)

	return WriteGenericMessage(writer, generic)
}
