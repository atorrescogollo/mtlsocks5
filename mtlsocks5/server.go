package mtlsocks5

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net"
	"os"
	"path/filepath"

	"tailscale.com/net/socks5"
)

type Server struct{}

func DefaultClientCACertPool(certDirectoryPath string) (*x509.CertPool, error) {
	pool := x509.NewCertPool()
	certFilePaths, err := filepath.Glob(filepath.Join(certDirectoryPath, "ca*.crt"))
	if err != nil {
		return nil, err
	}
	for _, p := range certFilePaths {
		certBytes, err := os.ReadFile(p)
		if err != nil {
			return nil, err
		}
		pool.AppendCertsFromPEM(certBytes)
	}
	return pool, nil
}

func DefaultServerTlsConfig(certPath string, keyPath string) (*tls.Config, error) {
	certificate, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, err
	}

	clientCAs, err := DefaultClientCACertPool("certs")
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		ClientCAs:    clientCAs,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{certificate},
		MinVersion:   tls.VersionTLS13,
		VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
			for _, rawCert := range rawCerts {
				cert, err := x509.ParseCertificate(rawCert)
				if err != nil {
					return err
				}
				//if err := cert.VerifyHostname("localhost"); err != nil {
				//	return err
				//}
				log.Printf("Verified certificate: %v", cert.DNSNames)
			}
			return nil
		},
	}, nil
}

func (s *Server) Serve(l net.Listener) error {
	tlsConfig, err := DefaultServerTlsConfig("certs/server.crt", "certs/server.key")
	if err != nil {
		return err
	}
	socksServer := socks5.Server{}
	return socksServer.Serve(
		tls.NewListener(l, tlsConfig),
	)
}
