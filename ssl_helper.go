package nebula_go_sdk

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"os"
)

// GetDefaultSSLConfig reads the files in the given path and returns a tls.Config object
func GetDefaultSSLConfig(rootCAPath, certPath, privateKeyPath string) (*tls.Config, error) {
	rootCA, err := openAndReadFile(rootCAPath)
	if err != nil {
		return nil, err
	}
	cert, err := openAndReadFile(certPath)
	if err != nil {
		return nil, err
	}
	privateKey, err := openAndReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}
	clientCert, err := tls.X509KeyPair(cert, privateKey)
	if err != nil {
		return nil, err
	}
	// parse root CA pem and add into CA pool
	// for self-signed cert, use the local cert as the root ca
	rootCAPool := x509.NewCertPool()
	ok := rootCAPool.AppendCertsFromPEM(rootCA)
	if !ok {
		return nil, fmt.Errorf("unable to append supplied cert into tls.Config, please make sure it is a valid certificate")
	}
	return &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      rootCAPool,
	}, nil
}

func openAndReadFile(path string) ([]byte, error) {
	// open file
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("unable to open file %s: %s", path, err)
	}
	// read file
	b, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("unable to ReadAll file %s: %s", path, err)
	}
	return b, nil
}
