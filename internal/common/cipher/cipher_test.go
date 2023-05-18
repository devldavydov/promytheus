package cipher

import (
	"bytes"
	"encoding/json"
	"io"
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

	var testMsg []byte
	for i := 0; i < 10000; i++ {
		testMsg = append(testMsg, []byte("Hello World!")...)
	}

	encrMsg, err := EncryptWithPublicKey(testMsg, pubKey)
	assert.NoError(t, err)

	decrMsg, err := DecryptWithPrivateKey(encrMsg, privKey)
	assert.NoError(t, err)

	assert.Equal(t, testMsg, decrMsg)
}

func TestEncDecrWithBuffer(t *testing.T) {
	type testStruct struct {
		A int    `json:"a"`
		B string `json:"b"`
	}
	testData := &testStruct{A: 123, B: "foobar"}

	privKey, pubKey, err := GenerateKeyPair(2048)
	assert.NoError(t, err)

	buf := NewEncBuffer(pubKey)
	json.NewEncoder(buf).Encode(testData)

	encData, err := io.ReadAll(buf)
	assert.NoError(t, err)
	encBuffer := bytes.NewBuffer(encData)

	type testEnvBufferClose struct {
		io.ReadWriter
		io.Closer
	}

	decRdr := NewDecReader(
		privKey, struct {
			io.ReadWriter
			io.Closer
		}{ReadWriter: encBuffer, Closer: nil},
	)
	decData := &testStruct{}
	json.NewDecoder(decRdr).Decode(decData)

	assert.Equal(t, testData, decData)
}
