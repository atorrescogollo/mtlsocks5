/*
Copyright © 2024 Álvaro Torres Cogollo <atorrescogollo@gmail.com>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// mgmtCmd represents the mgmt command
var mgmtCmd = &cobra.Command{
	Use:   "mgmt",
	Short: "Management commands",
	Long:  `Commands for managing certificates.`,
	//Run: func(cmd *cobra.Command, args []string) {
	//	fmt.Println("mgmt called")
	//},
}

func init() {
	rootCmd.AddCommand(mgmtCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mgmtCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// mgmtCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
