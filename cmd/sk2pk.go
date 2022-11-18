/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/sunliang711/crypto-km/utils"
)

// sk2pkCmd represents the sk2pk command
var sk2pkCmd = &cobra.Command{
	Use:   "sk2pk",
	Short: "private key to public key",
	Long:  `private key to public key`,
	Run: func(cmd *cobra.Command, args []string) {
		sk2pk(cmd, args)
	},
}

func init() {
	keyCmd.AddCommand(sk2pkCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sk2pkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sk2pkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	sk2pkCmd.Flags().String("sk", "", "secret key")
	sk2pkCmd.Flags().Bool("enter-sk", false, "enter secret key in safe way")
}

func sk2pk(cmd *cobra.Command, args []string) {
	var (
		sk  string
		err error
	)

	enterSk, _ := cmd.Flags().GetBool("enter-sk")

	if enterSk {
		sk, err = utils.ReadSecret(sk, "Enter secret key: ")
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}
	} else {
		sk, _ = cmd.Flags().GetString("sk")
		if sk == "" {
			fmt.Fprintf(os.Stderr, "No private key\n")
			os.Exit(1)
		}
	}

	pubkey, err := utils.PublieKey(sk)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("public key: %s\n", pubkey)

}
