package mtlsocks5

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"
	"net"

	"golang.org/x/net/proxy"
	"tailscale.com/net/socks5"
)

type Forwarder struct {
	ServerAddress string
}

func NewForwarder(serverAddress string) *Forwarder {
	return &Forwarder{
		ServerAddress: serverAddress,
	}
}

func DefaultForwarderTlsConfig(clientCertPath string, clientKeyPath string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(clientCertPath, clientKeyPath)
	if err != nil {
		return nil, err
	}

	clientCAs, err := DefaultClientCACertPool("certs")
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		RootCAs: clientCAs,
		VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
			for _, rawCert := range rawCerts {
				cert, err := x509.ParseCertificate(rawCert)
				if err != nil {
					return err
				}
				log.Printf("Verified certificate: %v\n", cert.DNSNames)
			}
			return nil
		},
		MinVersion:   tls.VersionTLS13,
		Certificates: []tls.Certificate{cert},
	}, nil
}

func forwarderDialer(serverAddress string) (proxy.Dialer, error) {
	// Create mTLS dialer to server
	tlsConfig, err := DefaultForwarderTlsConfig("certs/forwarder.crt", "certs/forwarder.key")
	if err != nil {
		return nil, err
	}
	tlsDialer := &tls.Dialer{
		Config: tlsConfig,
	}

	// Create SOCKS5 dialer through mTLS dialer to server
	socksDialer, err := proxy.SOCKS5("tcp", serverAddress, nil, tlsDialer)
	if err != nil {
		return nil, err
	}

	// Return the dialer
	return socksDialer, nil
}

func (f *Forwarder) Serve(l net.Listener) error {
	forwarderDialer, err := forwarderDialer(f.ServerAddress)
	if err != nil {
		return err
	}
	server := socks5.Server{
		Dialer: func(ctx context.Context, network, addr string) (net.Conn, error) {
			log.Printf("Dialing to %s\n", addr)
			return forwarderDialer.Dial(network, addr)
		},
	}
	return server.Serve(l)
}
