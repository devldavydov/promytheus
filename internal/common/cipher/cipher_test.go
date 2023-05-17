package cipher

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryptDecrypt(t *testing.T) {
	privKeyPath, pubKeyPath, err := generateTempKeyPair()
	assert.NoError(t, err)
	defer func() {
		os.Remove(privKeyPath)
		os.Remove(pubKeyPath)
	}()

	pubKey, err := PublicKeyFromFile(pubKeyPath)
	assert.NoError(t, err)

	privKey, err := PrivateKeyFromFile(privKeyPath)
	assert.NoError(t, err)

	testMsg := []byte("Hello world!")

	encrMsg, err := EncryptWithPublicKey(testMsg, pubKey)
	assert.NoError(t, err)

	decrMsg, err := DecryptWithPrivateKey(encrMsg, privKey)
	assert.NoError(t, err)

	assert.Equal(t, testMsg, decrMsg)
}

func generateTempKeyPair() (string, string, error) {
	privkey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", err
	}

	//
	privBytes := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privkey),
		},
	)

	pubASN1, err := x509.MarshalPKIXPublicKey(&privkey.PublicKey)
	if err != nil {
		return "", "", err
	}
	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubASN1,
	})

	//
	privFile, err := os.CreateTemp("", "key")
	if err != nil {
		return "", "", err
	}
	privFile.Write(privBytes)
	privFile.Close()

	pubFile, err := os.CreateTemp("", "key")
	if err != nil {
		return "", "", err
	}
	pubFile.Write(pubBytes)
	pubFile.Close()

	return privFile.Name(), pubFile.Name(), nil
}
