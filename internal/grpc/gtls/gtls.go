// Package for gRPC TLS utils.
package gtls

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"os"

	"google.golang.org/grpc/credentials"
)

type TLSServerSettings struct {
	ServerCertPath string
	ServerKeyPath  string
}

func NewOptionalTLSServerSettings(certPath, keyPath string) (*TLSServerSettings, error) {
	if certPath == "" && keyPath == "" {
		// No TLS
		return nil, nil
	}

	if (certPath == "" && keyPath != "") || (certPath != "" && keyPath == "") {
		return nil, errors.New("invalid TLS server settings")
	}

	return &TLSServerSettings{ServerCertPath: certPath, ServerKeyPath: keyPath}, nil
}

func (t *TLSServerSettings) Load() (credentials.TransportCredentials, error) {
	serverCert, err := tls.LoadX509KeyPair(t.ServerCertPath, t.ServerKeyPath)
	if err != nil {
		return nil, err
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert,
	}

	return credentials.NewTLS(config), nil
}

func LoadCACert(caCertPath string) (credentials.TransportCredentials, error) {
	if caCertPath == "" {
		// No TLS
		return nil, nil
	}

	pemServerCA, err := os.ReadFile(caCertPath)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("failed to add server CA's certificate")
	}

	config := &tls.Config{
		RootCAs: certPool,
	}

	return credentials.NewTLS(config), nil
}
