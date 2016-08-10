package main

import (
	"crypto/rsa"
	"flag"
	"fmt"

	"github.com/limoges/p2pnet/auth"
)

func main() {

	var priv *rsa.PrivateKey
	var hostkey []byte
	var err error
	var filepath string
	flag.Parse()

	filepath = flag.Arg(0)
	fmt.Printf("Reading private key from %v\n", filepath)
	if priv, err = auth.ReadPEMPrivateKey(filepath); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Getting hostkey as DER...")
	if hostkey, err = auth.GetPublicKeyAsDER(priv.PublicKey); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Hostkey length is %v\n", len(hostkey))
}
