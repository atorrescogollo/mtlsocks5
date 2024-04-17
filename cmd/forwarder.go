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

// forwarderCmd represents the forwarder command
var forwarderCmd = &cobra.Command{
	Use:   "forwarder",
	Short: "Starts a new mTLS SOCKS5 forwarder",
	Long:  `Starts a new mTLS SOCKS5 forwarder`,
	Run: func(cmd *cobra.Command, args []string) {
		listenAddr := cmd.Flag("listen-addr").Value.String()
		serverAddr := cmd.Flag("server-addr").Value.String()

		l, err := net.Listen("tcp", listenAddr)
		if err != nil {
			panic(err)
		}
		defer l.Close()

		log.Printf("Listening on %s", l.Addr().String())
		if err := mtlsocks5.NewForwarder(serverAddr).Serve(l); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(forwarderCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// forwarderCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// forwarderCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	forwarderCmd.Flags().StringP("listen-addr", "l", "0.0.0.0:1080", "Address to listen on")
	forwarderCmd.Flags().StringP("server-addr", "s", "", "Address of the mtlsocks5 server")
	forwarderCmd.MarkFlagRequired("server-addr")
}
