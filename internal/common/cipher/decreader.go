package cipher

import (
	"bytes"
	"crypto/rsa"
	"io"
)

// DecReader implements io.ReadCloser and decrypt input data.
type DecReader struct {
	io.ReadCloser

	privKey   *rsa.PrivateKey
	inp       io.Reader
	decBuf    bytes.Buffer
	decrypted bool
}

func NewDecReader(privKey *rsa.PrivateKey, inp io.Reader) *DecReader {
	return &DecReader{privKey: privKey, inp: inp}
}

func (d *DecReader) Read(p []byte) (n int, err error) {
	if !d.decrypted {
		encData, err := io.ReadAll(d.inp)
		if err != nil {
			return 0, err
		}

		decData, err := DecryptWithPrivateKey(encData, d.privKey)
		if err != nil {
			return 0, err
		}
		_, err = d.decBuf.Write(decData)
		if err != nil {
			return 0, err
		}

		d.decrypted = true
	}

	return d.decBuf.Read(p)
}

func (d *DecReader) Reset(inp io.Reader) {
	d.inp = inp
	d.decBuf.Reset()
	d.decrypted = false
}
