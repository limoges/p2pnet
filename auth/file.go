package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
)

// Read a PEM-formatted file into memory.
func ReadPEMPrivateKey(filepath string) (*rsa.PrivateKey, error) {

	var priv *rsa.PrivateKey
	var err error

	// Read the file.
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	// Decode the PEM-formatted data into a DER block.
	block, err := findPrivateKeyBlock(bytes)
	if err != nil {
		return nil, err
	}

	// Check if the PEM private key is password protected.
	if x509.IsEncryptedPEMBlock(block) {
		panic("Cannot handle encrypted keys for now")
	}

	if priv, err = x509.ParsePKCS1PrivateKey(block.Bytes); err != nil {
		fmt.Println(err)
		panic("Cannot parse private key")
	}

	///hostkey, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	///if err != nil {
	///	fmt.Println("Could not marshal public key into hostkey.")
	///	return nil, err
	///}

	///crypto = &Encryption{
	///	PrivateKey: privateKey,
	///	Hostkey:    hostkey,
	///}
	return priv, nil
}

func GetPublicKeyAsDER(pub *rsa.PublicKey) ([]byte, error) {
	return x509.MarshalPKIXPublicKey(pub)
}

func findPrivateKeyBlock(data []byte) (block *pem.Block, err error) {

	bytes := data
	for {
		block, rest := pem.Decode(bytes)
		if block == nil {
			return nil, ErrNoBlockFound
		}

		// Check if the block is a private key block.
		if block.Type == RSAPrivateKeyType {
			return block, nil
		}
		bytes = rest
	}
}

func EncodePEM(priv *rsa.PrivateKey) []byte {

	var block *pem.Block
	var derBytes []byte
	var pemBytes []byte

	derBytes = x509.MarshalPKCS1PrivateKey(priv)

	block = &pem.Block{
		Type:  RSAPrivateKeyType,
		Bytes: derBytes,
	}

	pemBytes = pem.EncodeToMemory(block)
	return pemBytes
}

func ParseHostkey(hostkey []byte) (publicKey *rsa.PublicKey, err error) {

	pub, err := x509.ParsePKIXPublicKey(hostkey)

	fmt.Println(pub)

	return nil, err
}
