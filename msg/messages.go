package msg

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
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

func SendReceive(conn net.Conn, message Message) (response Message, err error) {

	if err = Send(conn, message); err != nil {
		return nil, err
	}

	return Receive(conn)
}

func Receive(conn net.Conn) (Message, error) {

	var message Message
	var err error

	if message, err = Read(conn); err != nil {
		return message, err
	}

	fmt.Printf("%20v: RCV %25v from %20v\n",
		conn.LocalAddr(),
		Identifier(message.TypeId()),
		conn.RemoteAddr(),
	)
	return message, nil
}

// Read a message from the reader.
func Read(reader io.Reader) (Message, error) {

	var buf *bufio.Reader
	var generic GenericMessage
	var err error

	// Create a buffered reader because calls to read are blocking.
	// If we have an ill-formed message, we could be stuck trying to
	// wait for bytes that will never arrive.
	buf = bufio.NewReader(reader)

	if generic, err = ReadGenericMessage(buf); err != nil {
		return nil, err
	} else {
		return ConvertFromGeneric(generic)
	}
}

func Send(conn net.Conn, message Message) error {

	fmt.Printf("%20v: SND %25v to   %20v\n",
		conn.LocalAddr(),
		Identifier(message.TypeId()),
		conn.RemoteAddr(),
	)
	return Write(conn, message)
}

// Write a message to the writer.
func Write(writer io.Writer, m Message) error {

	var generic GenericMessage
	var err error

	if generic, err = ConvertToGeneric(m); err != nil {
		return err
	}

	return WriteGenericMessage(writer, generic)
}
