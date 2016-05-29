package auth

import (
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/limoges/p2pnet/cfg"
	"github.com/limoges/p2pnet/msg"
	"io/ioutil"
)

type Auth struct {
	HostKeyPath string
	block       *pem.Block
}

var (
	ErrNoBlockFound = errors.New("No block found in key file")
)

func New(conf *cfg.Configurations) (auth *Auth, err error) {

	auth = &Auth{}
	conf.Init(&auth.HostKeyPath, "", "HOSTKEY", "hostkey.pem")

	data, err := ioutil.ReadFile(auth.HostKeyPath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, ErrNoBlockFound
	}

	auth.block = block
	return auth, nil
}

func (a *Auth) handle(message msg.Message) {

	switch m := message.(type) {
	case msg.AuthSessionStart:
		a.StartSessionKeyEstablishment(m.HopHostKey)
	default:
		fmt.Printf("message is not handled by Auth:%v\n", message)
	}
}

func (a *Auth) createSessionID() uint32 {
	return 0
}

func (a *Auth) StartSessionKeyEstablishment(hopHostkey []byte) {

}
