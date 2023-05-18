// Package cipher provides RSA encryption/decryption.
package cipher

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

// PublicKeyFromFile loads public RSA key from file.
func PublicKeyFromFile(filePath string) (*rsa.PublicKey, error) {
	f, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return BytesToPublicKey(f)
}

// PrivateKeyFromFile loads private RSA key from file.
func PrivateKeyFromFile(filePath string) (*rsa.PrivateKey, error) {
	f, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return BytesToPrivateKey(f)
}

// GenerateKeyPair generates pair of RSA keys of given bits size.
func GenerateKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privkey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}
	return privkey, &privkey.PublicKey, nil
}

// GenerateKeyPair generates tmp files with pair of RSA keys of given bits size.
// Clients responsobility to move tmp files to specified destination.
func GenerateKeyPairFiles(bits int) (string, string, error) {
	privKey, pubKey, err := GenerateKeyPair(bits)
	if err != nil {
		return "", "", err
	}

	privFile, err := os.CreateTemp("", "key")
	if err != nil {
		return "", "", err
	}
	privFile.Write(PrivateKeyToBytes(privKey))
	privFile.Close()

	pubFile, err := os.CreateTemp("", "key")
	if err != nil {
		return "", "", err
	}
	pubBytes, err := PublicKeyToBytes(pubKey)
	if err != nil {
		return "", "", err
	}
	pubFile.Write(pubBytes)
	pubFile.Close()

	return privFile.Name(), pubFile.Name(), nil
}

// PrivateKeyToBytes converts RSA private key to slice of bytes.
func PrivateKeyToBytes(priv *rsa.PrivateKey) []byte {
	privBytes := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(priv),
		},
	)

	return privBytes
}

// PublicKeyToBytes converts RSA public key to slice of bytes.
func PublicKeyToBytes(pub *rsa.PublicKey) ([]byte, error) {
	pubASN1, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return nil, err
	}

	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubASN1,
	})

	return pubBytes, nil
}

// BytesToPrivateKey converts slice of bytes to RSA private key.
func BytesToPrivateKey(priv []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(priv)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			return nil, err
		}
	}
	key, err := x509.ParsePKCS1PrivateKey(b)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// BytesToPublicKey converts slice of bytes to RSA public key.
func BytesToPublicKey(pub []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pub)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			return nil, err
		}
	}
	ifc, err := x509.ParsePKIXPublicKey(b)
	if err != nil {
		return nil, err
	}
	key, ok := ifc.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not public key")
	}
	return key, nil
}

// EncryptWithPublicKey encrypts slice of bytes with RSA public key.
func EncryptWithPublicKey(msg []byte, pub *rsa.PublicKey) ([]byte, error) {
	hash := sha512.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, pub, msg, nil)
	if err != nil {
		return nil, err
	}
	return ciphertext, nil
}

// DecryptWithPrivateKey decrypts slice of bytes with RSA private key.
func DecryptWithPrivateKey(ciphertext []byte, priv *rsa.PrivateKey) ([]byte, error) {
	hash := sha512.New()
	plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, priv, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}
