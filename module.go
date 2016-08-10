package p2pnet

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/limoges/p2pnet/msg"
)

var (
	ErrModuleDoesNotHandle = errors.New("The present module cannot handle the message")
)

type Module interface {
	Name() string
	Addresses() (APIAddr, P2PAddr string)

	// Run should not return until the module has completed execution.
	Run() error

	// Handles messages that are appropriate for the module.
	Handle(source net.Conn, message msg.Message) error
}

func Run(m Module) error {

	// We launch the listeners, if they are supported by the module.
	apiAddr, p2pAddr := m.Addresses()

	if len(apiAddr) > 0 {
		fmt.Printf("%20v: %v: Listening API\n", m.Name(), apiAddr)
		if listener, err := net.Listen("tcp", apiAddr); err != nil {
			fmt.Printf("%v: Cannot bind on %v\n", m.Name(), apiAddr)
			return err
		} else {
			go listen(m, listener)
			defer listener.Close()
		}
	}

	if len(p2pAddr) > 0 {
		fmt.Printf("%20v: %v: Listening P2P\n", m.Name(), p2pAddr)
		if listener, err := net.Listen("tcp", p2pAddr); err != nil {
			fmt.Printf("%v: Cannot bind on %v\n", m.Name(), p2pAddr)
			return err
		} else {
			go listen(m, listener)
			defer listener.Close()
		}
	}

	// Then we call start the module's process.
	if err := m.Run(); err != nil {
		return err
	}
	return nil
}

func listen(m Module, ln net.Listener) {

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
		} else {
			go handle(m, conn)
		}
	}
}

func handle(m Module, conn net.Conn) {

	// fmt.Printf("%v: New connexion from %v.\n", m.Name(), conn.RemoteAddr())
	for {
		message, err := msg.Receive(conn)
		if err != nil {
			if err == io.EOF {
				// fmt.Printf("%v: Connexion with %v closed.\n", m.Name(), conn.RemoteAddr())
				return
			}
			log.Println(err)
			return
		}

		// fmt.Printf("%20v: %v: rcv %25v from %v\n",
		// 	m.Name(),
		// 	conn.LocalAddr(),
		// 	msg.Identifier(message.TypeId()),
		// 	conn.RemoteAddr(),
		// )

		if err := m.Handle(conn, message); err != nil {
			fmt.Println(err)
		}
	}
}
