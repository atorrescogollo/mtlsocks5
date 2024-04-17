/*
Copyright © 2024 Álvaro Torres Cogollo <atorrescogollo@gmail.com>
*/
package cmd

import (
	"log"
	"net"

	"github.com/atorrescogollo/mtlsocks5/mtlsocks5"
	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Starts a new mTLS SOCKS5 server",
	Long:  `Starts a new mTLS SOCKS5 server`,
	Run: func(cmd *cobra.Command, args []string) {
		listenAddr := cmd.Flag("listen-addr").Value.String()

		l, err := net.Listen("tcp", listenAddr)
		if err != nil {
			panic(err)
		}
		defer l.Close()

		log.Printf("Listening on %s", l.Addr().String())
		server := mtlsocks5.Server{}
		if err := server.Serve(l); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	serverCmd.Flags().StringP("listen-addr", "l", "0.0.0.0:1080", "Address to listen on")
}
