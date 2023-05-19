package cipher

import (
	"bytes"
	"crypto/rsa"
)

// EncBuffer implements io.ReadWriter with encryption.
// Write appends internal buffer witj plain data.
// Read returns encrypted data.
type EncBuffer struct {
	pubKey  *rsa.PublicKey
	buf     bytes.Buffer
	encBuf  bytes.Buffer
	encoded bool
}

func NewEncBuffer(pubKey *rsa.PublicKey) *EncBuffer {
	return &EncBuffer{pubKey: pubKey}
}

func (e *EncBuffer) Write(p []byte) (n int, err error) {
	return e.buf.Write(p)
}

func (e *EncBuffer) Read(p []byte) (n int, err error) {
	if !e.encoded {
		encData, err := EncryptWithPublicKey(e.buf.Bytes(), e.pubKey)
		if err != nil {
			return 0, err
		}

		_, err = e.encBuf.Write(encData)
		if err != nil {
			return 0, err
		}
		e.encoded = true
	}

	return e.encBuf.Read(p)
}
