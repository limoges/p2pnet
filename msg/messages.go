package msg

import (
	"errors"
	"fmt"
	"io"
)

const (
	HeaderLength = 4
	IPLength     = 16
)

var (
	ErrDataTooShort = errors.New("unmarshal: data is shorter than expected")
)

func createMessage(messageType int, payloadBuf []byte) (data []byte) {

	panic("Remove this function")
	return data
}

type Message interface {
	TypeId() uint16
}

// Read a message from the reader.
func ReadMessage(reader io.Reader) (Message, error) {

	if generic, err := ReadGenericMessage(reader); err != nil {
		return nil, err
	} else {
		fmt.Println("Reading:", generic)
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
	fmt.Println("Writing:", generic)

	return WriteGenericMessage(writer, generic)
}
