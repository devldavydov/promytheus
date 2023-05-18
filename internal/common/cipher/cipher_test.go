package cipher

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryptDecrypt(t *testing.T) {
	privKeyPath, pubKeyPath, err := GenerateKeyPairFiles(2048)
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
