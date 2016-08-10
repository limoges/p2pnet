package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"errors"
	"hash"
	"io"
)

const (
	RSAPrivateKeyType                = "RSA PRIVATE KEY"
	RSAPublicKeyType                 = "RSA PUBLIC KEY"
	DefaultAsymmetricKeyLengthInBits = 4096
	DefaultSymmetricKeyLengthInBytes = 16
)

type Encryption struct {
	Hostkey    []byte
	PrivateKey *rsa.PrivateKey
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

func CheckMAC(message, messageMAC, key []byte) bool {
	var mac hash.Hash
	var expectedMAC []byte

	mac = hmac.New(sha256.New, key)
	mac.Write(message)
	expectedMAC = mac.Sum(nil)
	return hmac.Equal(messageMAC, expectedMAC)
}

func ComputeMAC(message, key []byte) []byte {

	var mac hash.Hash

	mac = hmac.New(sha256.New, key)
	mac.Write(message)
	return mac.Sum(nil)
}

func EncryptAES(key, plaintext []byte) ([]byte, error) {

	var block cipher.Block
	var ciphertext []byte
	var iv []byte
	var err error

	if block, err = aes.NewCipher(key); err != nil {
		return nil, err
	}

	// The IV needs to be unique. It's common to include it at the beginning
	// of the ciphertext.
	ciphertext = make([]byte, aes.BlockSize+len(plaintext))

	iv = ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return ciphertext, nil
}

func DecryptAES(key, ciphertext []byte) ([]byte, error) {

	var block cipher.Block
	var err error
	var iv []byte

	if block, err = aes.NewCipher(key); err != nil {
		return nil, err
	}

	iv = ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext, nil
}

func EncryptPKCS(pub *rsa.PublicKey, msg []byte) ([]byte, error) {

	return rsa.EncryptPKCS1v15(rand.Reader, pub, msg)
}

func DecryptPKCS(priv *rsa.PrivateKey, ciphertext []byte) ([]byte, error) {

	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}

func EncryptAESWithHMAC(plaintext, secret, hmac []byte) ([]byte, error) {

	var ciphertext, signature, encrypted []byte
	var err error

	if ciphertext, err = EncryptAES(secret, plaintext); err != nil {
		return nil, err
	}

	signature = ComputeMAC(ciphertext, hmac)

	// Validate the length of this signature.
	if len(signature) != 32 {
		panic("Signature is expected to be 32 bytes")
	}

	encrypted = make([]byte, 0, len(ciphertext)+len(signature))
	encrypted = append(encrypted, signature...)
	encrypted = append(encrypted, ciphertext...)

	return encrypted, nil
}

func DecryptAESWithHMAC(encrypted, secret, hmac []byte) ([]byte, error) {

	var ciphertext, signature []byte
	var valid bool

	signature = encrypted[:32]
	ciphertext = make([]byte, len(encrypted[32:]))
	copy(ciphertext, encrypted[32:])

	if valid = CheckMAC(ciphertext, signature, hmac); !valid {
		return nil, errors.New("Signature does not match message.")
	}

	return DecryptAES(secret, ciphertext)
}
