package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
)

func GeneratePrivatePubKey() (pubKey []byte, privateKey []byte, err error) {
	pubKey, privateKey, err = ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	return
}

func Sign(privateKey []byte, message string) []byte {
	return ed25519.Sign(privateKey, []byte(message))
}

func Verify(pubKey []byte, message, signature []byte) bool {
	return ed25519.Verify(pubKey, message, signature)
}
