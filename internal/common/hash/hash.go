package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func HmacSHA256(data, key string) string {
	hmac := hmac.New(sha256.New, []byte(key))
	hmac.Write([]byte(data))
	return hex.EncodeToString(hmac.Sum(nil))
}

func HmacEqual(hmac1, hmac2 string) bool {
	hmac1b, _ := hex.DecodeString(hmac1)
	hmac2b, _ := hex.DecodeString(hmac2)
	return hmac.Equal(hmac1b, hmac2b)
}
