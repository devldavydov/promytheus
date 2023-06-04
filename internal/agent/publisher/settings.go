package publisher

import (
	"crypto/rsa"
	"net"
	"time"

	"google.golang.org/grpc/credentials"
)

const (
	_defaultRequestTimeout  = 15 * time.Second
	_defaultShutdownTimeout = 5 * time.Second
)

// EncryptionSettings - settings for publishers encryption
type EncryptionSettings struct {
	CryptoPubKey   *rsa.PublicKey
	TLSCredentials credentials.TransportCredentials
}

type PublisherExtraSettings struct {
	HmacKey         *string
	EncrSettings    EncryptionSettings
	ShutdownTimeout *time.Duration
	HostIP          net.IP
}
