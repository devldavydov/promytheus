// Package hash providies functions to work with HMAC.
package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// HmacSHA256 get hash value for data and key.
func HmacSHA256(data, key string) string {
	hmac := hmac.New(sha256.New, []byte(key))
	hmac.Write([]byte(data))
	return hex.EncodeToString(hmac.Sum(nil))
}

// HmacEqual compares to hmac.
func HmacEqual(hmac1, hmac2 string) bool {
	hmac1b, _ := hex.DecodeString(hmac1)
	hmac2b, _ := hex.DecodeString(hmac2)
	return hmac.Equal(hmac1b, hmac2b)
}
