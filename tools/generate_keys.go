package main

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"

	"github.com/limoges/p2pnet/auth"
)

func main() {

	fmt.Println("This utility will generate a public-private key pair.")
	fmt.Println("The key pair is referred as hostkey, is generated using RSA")
	fmt.Println("and is 4096 bits long. It will be stored on disk in PEM format.")

	var priv *rsa.PrivateKey
	var data []byte
	var hostkey []byte
	var err error

	fmt.Println("Generating new public-private key pair...")
	if priv, err = auth.GenerateKey(); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Encoding keys to PEM...")
	data = auth.EncodePEM(priv)

	fmt.Println("Getting hostkey as DER...")
	if hostkey, err = auth.GetPublicKeyAsDER(priv.PublicKey); err != nil {
		panic(err)
	}
	fmt.Printf("Hostkey length is %v\n", len(hostkey))

	fmt.Println("Writing keys to file...")
	if err = ioutil.WriteFile("keys.pem", data, 0777); err != nil {
		fmt.Println(err)
		return
	}
}
