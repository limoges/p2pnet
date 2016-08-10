package auth

import (
	"crypto/rsa"
	"errors"
	"math/rand"
	"time"

	"github.com/limoges/p2pnet/msg"
)

type Session struct {
	Id              uint32
	SharedKey       []byte
	LocalHMAC       []byte
	RemoteHMAC      []byte
	RemotePublicKey *rsa.PublicKey
}

func (a *Auth) localUnusedSessionId() (uint32, error) {

	const MaximumSessionIdAttempts = 100
	var id uint32
	var alreadyInUse bool
	var count int

	// With the estimated number of simultaneous connexion very low,
	// it is suggested to use a 32 bits session ID which will have
	// about odds of a collision in about 1 in 10 millions for 30 or
	// simultaneous sessions.
	// Therefore, using the time as a random number generator seed seems
	// perfectly fine.

	seed := time.Now().UnixNano()
	source := rand.NewSource(seed)
	rng := rand.New(source)

	// We generate a random number and then try again until we can find
	// a number that is not currently in use. Again, there shouldn't be too
	// many collisions locally.
	id = rng.Uint32()
	count = 1
	for {
		if count == MaximumSessionIdAttempts {
			return 0, errors.New("Could not generate a new session id")
		}
		if _, alreadyInUse = a.Sessions[id]; !alreadyInUse {
			break
		}
		id = rng.Uint32()
		count = count + 1
	}
	return id, nil
}

func NewIncomingSession(a *Auth, encryptedKey, encryptedHMAC []byte) (*Session, error) {

	var id uint32
	var localHMAC []byte
	var session *Session
	var remoteHMAC []byte
	var sharedKey []byte
	var err error

	if sharedKey, err = DecryptPKCS(a.PrivateKey, encryptedKey); err != nil {
		return nil, err
	}

	if remoteHMAC, err = DecryptPKCS(a.PrivateKey, encryptedHMAC); err != nil {
		return nil, err
	}

	if localHMAC, err = GenerateNewSymmetricKey(); err != nil {
		return nil, err
	}

	if id, err = a.localUnusedSessionId(); err != nil {
		return nil, err
	}

	session = &Session{
		Id:         id,
		SharedKey:  sharedKey,
		LocalHMAC:  localHMAC,
		RemoteHMAC: remoteHMAC,
	}
	return session, nil
}
func NewSession(a *Auth) (*Session, error) {

	var id uint32
	var sharedKey []byte
	var localHMAC []byte
	var session *Session
	var err error

	if id, err = a.localUnusedSessionId(); err != nil {
		return nil, err
	}

	if sharedKey, err = GenerateNewSymmetricKey(); err != nil {
		return nil, err
	}

	if localHMAC, err = GenerateNewSymmetricKey(); err != nil {
		return nil, err
	}

	session = &Session{
		Id:        id,
		SharedKey: sharedKey,
		LocalHMAC: localHMAC,
	}
	return session, nil
}
func (s *Session) DecryptRemoteHMAC(priv *rsa.PrivateKey, encryptedHMAC []byte) error {

	var err error
	if s.RemoteHMAC, err = DecryptPKCS(priv, encryptedHMAC); err != nil {
		return err
	}

	return nil
}

func (s *Session) Validate() error {
	return nil
}

func (s *Session) CreateHandshake1(pub *rsa.PublicKey) (*msg.AuthSessionHS1, error) {

	var encryptedKey []byte
	var encryptedHMAC []byte
	var handshake1 msg.AuthHandshake1
	var payload []byte
	var session1 msg.AuthSessionHS1
	var err error

	if encryptedKey, err = EncryptPKCS(pub, s.SharedKey); err != nil {
		return nil, err
	}

	if encryptedHMAC, err = EncryptPKCS(pub, s.LocalHMAC); err != nil {
		return nil, err
	}

	if len(encryptedKey) != 512 {
		panic("Problem with the algorithm")
	}

	if len(encryptedHMAC) != 512 {
		panic("Problem with the algorithm")
	}

	handshake1 = msg.AuthHandshake1{}
	copy(handshake1.EncryptedKey[:], encryptedKey)
	copy(handshake1.EncryptedHMAC[:], encryptedHMAC)

	if payload, err = buildPayload(handshake1); err != nil {
		return nil, err
	}

	session1 = msg.AuthSessionHS1{
		SessionId:        s.Id,
		HandshakePayload: payload,
	}

	return &session1, nil
}

func (s *Session) CreateHandshake2(pub *rsa.PublicKey) (*msg.AuthSessionHS2, error) {

	var encryptedHMAC []byte
	var err error
	var handshake2 msg.AuthHandshake2
	var payload []byte
	var session2 msg.AuthSessionHS2

	if encryptedHMAC, err = EncryptPKCS(pub, s.LocalHMAC); err != nil {
		return nil, err
	}

	if len(encryptedHMAC) != 512 {
		panic("Problem with algoritm")
	}
	handshake2 = msg.AuthHandshake2{}
	copy(handshake2.EncryptedHMAC[:], encryptedHMAC)

	if payload, err = buildPayload(handshake2); err != nil {
		return nil, err
	}

	session2 = msg.AuthSessionHS2{
		SessionId:        s.Id,
		HandshakePayload: payload,
	}
	return &session2, nil
}

func (s *Session) Encrypt(plaintext []byte) ([]byte, error) {
	return EncryptAESWithHMAC(plaintext, s.SharedKey, s.LocalHMAC)
}

func (s *Session) Decrypt(ciphertext []byte) ([]byte, error) {
	return DecryptAESWithHMAC(ciphertext, s.SharedKey, s.RemoteHMAC)
}
