/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"crypto/ecdsa"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
	"github.com/sunliang711/crypto-km/utils"
)

// pk2addressCmd represents the pk2address command
var pk2addressCmd = &cobra.Command{
	Use:   "pk2address",
	Short: "public key to address",
	Long:  `public key to address by blockchain type`,
	Run: func(cmd *cobra.Command, args []string) {
		pk2address(cmd, args)
	},
}

func init() {
	keyCmd.AddCommand(pk2addressCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pk2addressCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pk2addressCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	pk2addressCmd.Flags().StringP("pubkey", "p", "", "public key")
	pk2addressCmd.Flags().StringP("type", "t", "eth", "blockchain type, eg: btc eth sol ada dot ksm")
}

func pk2address(cmd *cobra.Command, args []string) {
	pubkey, _ := cmd.Flags().GetString("pubkey")
	if pubkey == "" {
		fmt.Fprintf(os.Stderr, "no pubkey\n")
		os.Exit(1)
	}

	blockchainType, _ := cmd.Flags().GetString("type")
	switch blockchainType {
	case "btc":
	case "eth":
		address, err := pk2ethAddress(pubkey)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		fmt.Printf("address: %s\n", address)
	case "sol":

	default:
		fmt.Fprintf(os.Stderr, "blockchain %s not supported\n", blockchainType)
		os.Exit(1)

	}
}

func pk2ethAddress(pubkey string) (string, error) {
	var (
		publicKey *ecdsa.PublicKey
	)
	p, err := hexutil.Decode(pubkey)
	if err != nil {
		return "", err
	}
	compressed, err := utils.IsCompressedPublicKey(p)
	if err != nil {
		return "", err
	}

	if compressed {
		fmt.Printf("compressed\n")
		publicKey, err = utils.RecoverPublicKeyFromCompressed(p)
		if err != nil {
			return "", err
		}
	} else {
		fmt.Printf("uncompressed\n")
		publicKey, err = crypto.UnmarshalPubkey(p)
		if err != nil {
			return "", err
		}
	}
	fmt.Printf("publickey: %+v\n", publicKey)
	address := crypto.PubkeyToAddress(*publicKey)
	return address.Hex(), nil

}
