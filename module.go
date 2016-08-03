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

	fmt.Printf("Starting module '%v'\n", m.Name())
	// We launch the listeners, if they are supported by the module.
	apiAddr, p2pAddr := m.Addresses()

	if len(apiAddr) > 0 {
		fmt.Printf("Launching API on %v\n", apiAddr)
		if listener, err := net.Listen("tcp", apiAddr); err != nil {
			fmt.Printf("Cannot bind on %v\n", apiAddr)
			return err
		} else {
			go listen(m, listener)
			defer listener.Close()
		}
	}

	if len(p2pAddr) > 0 {
		fmt.Printf("Launching P2P on %v\n", p2pAddr)
		if listener, err := net.Listen("tcp", p2pAddr); err != nil {
			fmt.Printf("Cannot bind on %v\n", p2pAddr)
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

	log.Printf("New connexion from %v.\n", conn.RemoteAddr())
	for {
		message, err := msg.ReadMessage(conn)
		if err != nil {
			if err == io.EOF {
				log.Printf("Connexion with %v closed.\n", conn.RemoteAddr())
				return
			}
			log.Println(err)
			return
		}

		if err := m.Handle(conn, message); err != nil {
			fmt.Println(err)
		}
	}
}
