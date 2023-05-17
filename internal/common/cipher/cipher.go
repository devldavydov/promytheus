package cipher

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
)

func LoadPublicKey(filePath string) (*rsa.PublicKey, error) {
	return nil, nil
}

func LoadPrivateKey(filePath string) (*rsa.PrivateKey, error) {
	return nil, nil
}

func EncryptWithPublicKey(msg []byte, pub *rsa.PublicKey) ([]byte, error) {
	hash := sha512.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, pub, msg, nil)
	if err != nil {
		return nil, err
	}
	return ciphertext, nil
}

func DecryptWithPrivateKey(ciphertext []byte, priv *rsa.PrivateKey) ([]byte, error) {
	hash := sha512.New()
	plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, priv, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}
