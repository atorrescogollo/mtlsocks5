/*
Copyright © 2024 Álvaro Torres Cogollo <atorrescogollo@gmail.com>
*/
package cmd

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/atorrescogollo/mtlsocks5/utils"
	"github.com/spf13/cobra"
)

// newCaCmd represents the newCa command
var newCaCmd = &cobra.Command{
	Use:   "new-ca",
	Short: "Creates a new CA certificate and key pair",
	Long:  `Creates a new CA certificate and key pair. The certificate will be saved in the certs directory by default.`,
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat("certs/ca.crt"); err == nil {
			log.Panicf("CA certificate already exists in certs/ca.crt")
		}
		if _, err := os.Stat("certs/ca.key"); err == nil {
			log.Panicf("CA private key already exists in certs/ca.key")
		}

		log.Println("Creating new CA private key")
		priv, err := rsa.GenerateKey(rand.Reader, 4096)
		if err != nil {
			panic(err)
		}
		log.Println("Creating new CA certificate")
		cert, err := utils.NewCACertificate(x509.Certificate{
			SerialNumber: big.NewInt(1),
			Subject: pkix.Name{
				CommonName: "mtlsocks5 CA",
			},
			NotBefore:             time.Now(),
			NotAfter:              time.Now().AddDate(10, 0, 0),
			BasicConstraintsValid: true,
			IsCA:                  true,
			KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		}, *priv)
		if err != nil {
			panic(err)
		}

		log.Println("Saving CA certificate and key")
		if err := os.MkdirAll("certs", 0o700); err != nil {
			panic(err)
		}
		if err := utils.SaveCertificateAsPEM(*cert, "certs/ca.crt"); err != nil {
			panic(err)
		}
		if err := utils.SavePrivateKeyAsPEM(*priv, "certs/ca.key"); err != nil {
			panic(err)
		}
	},
}

func init() {
	mgmtCmd.AddCommand(newCaCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCaCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newCaCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
