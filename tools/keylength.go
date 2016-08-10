package main

import (
	"crypto/rsa"
	"encoding/hex"
	"fmt"

	"github.com/limoges/p2pnet/auth"
)

func main() {

	var secret []byte
	var hmac []byte
	var plaintext []byte
	var encrypted, decrypted []byte
	var err error

	plaintext = []byte("A B C D E F G H I J K L M N O P Q R S T U V W X Y Z")

	secret = asymmetric()
	hmac = generateHMAC()

	if encrypted, err = auth.EncryptAESWithHMAC(plaintext, secret, hmac); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Encrypted length is %v\n", len(encrypted))

	if decrypted, err = auth.DecryptAESWithHMAC(encrypted, secret, hmac); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Encrypted Message:%v\n", string(encrypted))
	fmt.Printf("Original Message :%v\n", string(plaintext))
	fmt.Printf("Decrypted Message:%v\n", string(decrypted))
}

func generateHMAC() []byte {

	var hmac []byte
	var err error

	if hmac, err = auth.GenerateNewSymmetricKey(); err != nil {
		panic(err)
	}

	return hmac
}

func asymmetric() []byte {
	var private *rsa.PrivateKey
	var key []byte
	var encrypted []byte
	var decrypted []byte
	var err error

	fmt.Println("Generating RSA private key.")
	if private, err = auth.GenerateKey(); err != nil {
		panic(err)
	}

	if key, err = auth.GenerateNewSymmetricKey(); err != nil {
		panic(err)
	}
	fmt.Printf("Generated %v-bytes symmetric key.\n", len(key))

	if encrypted, err = auth.EncryptPKCS(&private.PublicKey, key); err != nil {
		panic(err)
	}
	fmt.Printf("Encrypted symmetric key. Encrypted length %v.\n", len(encrypted))

	if decrypted, err = auth.DecryptPKCS(private, encrypted); err != nil {
		panic(err)
	}

	fmt.Printf("Unencrypted key:%v\n", hex.EncodeToString(key))
	fmt.Printf("Decrypted key  :%v\n", hex.EncodeToString(decrypted))
	fmt.Printf("Encrypted key  :%v\n", hex.EncodeToString(encrypted))
	return key
}
