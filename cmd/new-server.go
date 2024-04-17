/*
Copyright © 2024 Álvaro Torres Cogollo <atorrescogollo@gmail.com>
*/
package cmd

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/atorrescogollo/mtlsocks5/utils"
	"github.com/spf13/cobra"
)

// newServerCmd represents the newServer command
var newServerCmd = &cobra.Command{
	Use:   "new-server",
	Short: "Creates a new server certificate and key pair",
	Long:  `Creates a new server certificate and key pair. The certificate will be saved in the certs directory by default.`,
	Run: func(cmd *cobra.Command, args []string) {
		serverName := cmd.Flag("server-name").Value.String()
		certPath := fmt.Sprintf("certs/%s.crt", serverName)
		keyPath := fmt.Sprintf("certs/%s.key", serverName)

		if _, err := os.Stat(certPath); err == nil {
			log.Panicf("Server certificate already exists in %s", certPath)
		}
		if _, err := os.Stat(keyPath); err == nil {
			log.Panicf("Server private key already exists in %s", keyPath)
		}

		log.Println("Loading CA certificate and key")
		caPrivKey, err := utils.LoadPEMPrivateKeyFromPath("certs/ca.key")
		if err != nil {
			panic(err)
		}
		caCert, err := utils.LoadPEMCertificateFromPath("certs/ca.crt")
		if err != nil {
			panic(err)
		}

		log.Println("Creating new server private key")
		priv, err := rsa.GenerateKey(rand.Reader, 4096)
		if err != nil {
			panic(err)
		}
		log.Printf("Creating new server certificate for %s", serverName)
		cert, err := utils.NewPeerCertificate(x509.Certificate{
			SerialNumber: big.NewInt(1),
			Subject: pkix.Name{
				CommonName: serverName,
			},
			DNSNames:    []string{serverName},
			NotBefore:   time.Now(),
			NotAfter:    time.Now().AddDate(1, 0, 0),
			ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
			KeyUsage:    x509.KeyUsageDigitalSignature,
		}, *priv, *caCert, *caPrivKey)
		if err != nil {
			panic(err)
		}

		log.Printf("Saving server certificate and key in %s and %s", certPath, keyPath)
		if err := utils.SaveCertificateAsPEM(*cert, certPath); err != nil {
			panic(err)
		}
		if err := utils.SavePrivateKeyAsPEM(*priv, keyPath); err != nil {
			panic(err)
		}
	},
}

func init() {
	mgmtCmd.AddCommand(newServerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newServerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newServerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	newServerCmd.Flags().StringP("server-name", "n", "", "Common name for the server certificate")
	newServerCmd.MarkFlagRequired("server-name")
}
