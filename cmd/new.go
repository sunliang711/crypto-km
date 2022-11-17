/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/sunliang711/crypto-km/utils"
	"github.com/tyler-smith/go-bip39"
)

// ckm bip39
//     new [--words 12] [-o output]

const ENTROPY_BIT_SIZE_12 = 32 * 4
const ENTROPY_BIT_SIZE_15 = 32 * 5
const ENTROPY_BIT_SIZE_18 = 32 * 6
const ENTROPY_BIT_SIZE_21 = 32 * 8
const ENTROPY_BIT_SIZE_24 = 32 * 8

// [128, 256]
// 32 * 4 = 12 words
// 32 * 5 = 15 words
// 32 * 6 = 18 words
// 32 * 7 = 21 words
// 32 * 8 = 24 words

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "create a random mnemonic",
	Long:  `create a random mnemonic`,
	Run: func(cmd *cobra.Command, args []string) {
		createMnemonic(cmd, args)
	},
}

func init() {
	bip39Cmd.AddCommand(newCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	newCmd.Flags().Uint8("words", 12, "mnemonic words count")
	newCmd.Flags().StringP("output", "o", "", "output file, print to stdout if empty")
}

func createMnemonic(cmd *cobra.Command, args []string) {

	words, _ := cmd.Flags().GetUint8("words")
	output, _ := cmd.Flags().GetString("output")
	bitSize := 0

	switch words {
	case 12:
		bitSize = ENTROPY_BIT_SIZE_12
	case 15:
		bitSize = ENTROPY_BIT_SIZE_15
	case 18:
		bitSize = ENTROPY_BIT_SIZE_18
	case 21:
		bitSize = ENTROPY_BIT_SIZE_21
	case 24:
		bitSize = ENTROPY_BIT_SIZE_24
	default:
		fmt.Printf("invalid words,only support 12, 15, 18, 21, 24")
		return
	}

	entropy, err := bip39.NewEntropy(bitSize)
	if err != nil {
		fmt.Fprintf(os.Stderr, "new entropy failed: %s\n", err)
		return
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		fmt.Fprintf(os.Stderr, "new mnemonic failed: %s\n", err)
		return
	}

	if output != "" {
		err = utils.WriteFileWhenNotExists(output, []byte(mnemonic), 0600)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
		}
	} else {
		fmt.Printf("%s\n", mnemonic)
	}
}
