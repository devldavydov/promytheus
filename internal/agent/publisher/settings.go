package publisher

import (
	"crypto/rsa"
	"net"
	"time"
)

const (
	_defaultRequestTimeout  = 15 * time.Second
	_defaultShutdownTimeout = 5 * time.Second
)

// EncryptionSettings - settings for publishers encryption
type EncryptionSettings struct {
	CryptoPubKey *rsa.PublicKey
}

type PublisherExtraSettings struct {
	HmacKey         *string
	EncrSettings    EncryptionSettings
	ShutdownTimeout *time.Duration
	HostIP          net.IP
}
