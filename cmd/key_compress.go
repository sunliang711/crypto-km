/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
)

// compressCmd represents the compress command
var compressCmd = &cobra.Command{
	Use:   "compress",
	Short: "compress or decompress public key",
	Long:  `compress or decompress public key`,
	Run: func(cmd *cobra.Command, args []string) {
		compress(cmd, args)
	},
}

func init() {
	keyCmd.AddCommand(compressCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// compressCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// compressCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	compressCmd.Flags().BoolP("decompress", "d", false, "decompress or not")
	compressCmd.Flags().String("pubkey", "", "public key(eg: 0x...)")
}

func compress(cmd *cobra.Command, args []string) {
	decompress, _ := cmd.Flags().GetBool("decompress")
	pubkey, _ := cmd.Flags().GetString("pubkey")
	if pubkey == "" {
		fmt.Fprintf(os.Stderr, "no pubkey\n")
		os.Exit(1)
	}

	pubkeyBytes, err := hexutil.Decode(pubkey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "decode public key error: %s\n", err)
		os.Exit(1)
	}

	if decompress {
		p, err := crypto.DecompressPubkey(pubkeyBytes)
		if err != nil {
			fmt.Fprintf(os.Stderr, "decompress pubkey error: %s\n", err)
			os.Exit(1)
		}

		decompressedBytes := crypto.FromECDSAPub(p)
		fmt.Printf("decompressed: %s\n", hexutil.Encode(decompressedBytes))
	} else {
		p, err := crypto.UnmarshalPubkey(pubkeyBytes)
		if err != nil {
			fmt.Fprintf(os.Stderr, "unmarshal pubkey error: %s\n", err)
			os.Exit(1)
		}

		compressed := crypto.CompressPubkey(p)
		fmt.Printf("compressed: %s\n", hexutil.Encode(compressed))
	}

}
