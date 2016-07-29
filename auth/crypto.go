package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
)

const (
	RSAPrivateKeyType                = "RSA PRIVATE KEY"
	RSAPublicKeyType                 = "RSA PUBLIC KEY"
	DefaultAsymmetricKeyLengthInBits = 4096
	DefaultSymmetricKeyLengthInBytes = 64
)

type Encryption struct {
	Hostkey    []byte
	PrivateKey *rsa.PrivateKey
}

// Read a PEM-formatted file into memory.
func ReadKeys(filepath string) (crypto *Encryption, err error) {

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

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		fmt.Println(err)
		panic("Cannot parse private key")
		return nil, err
	}

	hostkey, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		fmt.Println("Could not marshal public key into hostkey.")
		return nil, err
	}

	crypto = &Encryption{
		PrivateKey: privateKey,
		Hostkey:    hostkey,
	}
	return crypto, nil
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

func GenerateKey() (*rsa.PrivateKey, error) {

	var priv *rsa.PrivateKey
	var err error

	if priv, err = rsa.GenerateKey(rand.Reader, DefaultAsymmetricKeyLengthInBits); err != nil {
		return nil, err
	}

	priv.Precompute()

	if err = priv.Validate(); err != nil {
		return nil, err
	}
	return priv, nil
}

func MarshalPublicKey(pub *rsa.PublicKey) ([]byte, error) {

	var derBytes []byte
	var err error

	derBytes, err = x509.MarshalPKIXPublicKey(pub)
	return derBytes, err
}

func ParsePublicKey(derBytes []byte) (*rsa.PublicKey, error) {

	var pub *rsa.PublicKey
	var err error
	var publicKey interface{}
	var properFormat bool

	// Parse the public key into the rsa.PublicKey
	if publicKey, err = x509.ParsePKIXPublicKey(derBytes); err != nil {
		return nil, err
	}

	// Check that the key is in the proper format
	if pub, properFormat = publicKey.(*rsa.PublicKey); !properFormat {
		return nil, errors.New("The pubublic key is not RSA.")
	}

	return pub, nil
}

func GenerateNewSymmetricKey() ([]byte, error) {

	var key []byte
	var err error

	key = make([]byte, DefaultSymmetricKeyLengthInBytes)

	if _, err = rand.Read(key); err != nil {
		return nil, err
	}
	return key, nil
}

func AsymmetricEncrypt(pub *rsa.PublicKey, msg []byte) ([]byte, error) {

	var ciphertext []byte
	var err error

	ciphertext, err = rsa.EncryptPKCS1v15(rand.Reader, pub, msg)
	return ciphertext, err
}

func AsymmetricDecrypt(priv *rsa.PrivateKey, ciphertext []byte) ([]byte, error) {

	var out []byte
	var err error

	out, err = rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
	return out, err
}
