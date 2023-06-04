package middleware

import (
	"crypto/rsa"
	"net/http"

	"github.com/devldavydov/promytheus/internal/common/cipher"
)

// Decrypt is a RSA decryption middleware.
type Decrypt struct {
	privKey *rsa.PrivateKey
}

func NewDecrpyt(privKey *rsa.PrivateKey) *Decrypt {
	return &Decrypt{privKey: privKey}
}

func (d *Decrypt) Handle(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if d.privKey != nil {
			r.Body = cipher.NewDecReader(d.privKey, r.Body)
		}

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
