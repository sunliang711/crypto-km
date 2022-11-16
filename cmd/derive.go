/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/spf13/cobra"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

const QUOTE_PREFIX = 0x80000000

// deriveCmd represents the derive command
var deriveCmd = &cobra.Command{
	Use:   "derive",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		derive(cmd, args)
	},
}

func init() {
	bip39Cmd.AddCommand(deriveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deriveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deriveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	deriveCmd.Flags().String("mnemonic", "", "mnemonic")
	deriveCmd.Flags().String("password", "", "password")
	deriveCmd.Flags().String("path", "", "hd path")
	deriveCmd.Flags().StringP("output", "o", "", "output file, print to stdout if empty")
	deriveCmd.Flags().Bool("with-seed", false, "show seed")
}

func derive(cmd *cobra.Command, args []string) {
	var (
		outputContent string
		sk            string
		pubkey        string
	)

	output, _ := cmd.Flags().GetString("output")
	// TODO: read mnemonic with no echo
	mnemonic, _ := cmd.Flags().GetString("mnemonic")
	// TODO: read password with no echo
	password, _ := cmd.Flags().GetString("password")

	path, _ := cmd.Flags().GetString("path")
	withSeed, _ := cmd.Flags().GetBool("with-seed")

	seed := bip39.NewSeed(mnemonic, password)

	rootKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Create master key error: %s\n", err)
		return
	}

	// return root sk when path is empty
	if path == "" {
		sk = hexutil.Encode(rootKey.Key)
		pubkey = hexutil.Encode(rootKey.PublicKey().Key)
	} else {
		derivedKey, err := _derive(rootKey, path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "derive error: %s\n", err)
			return
		}
		sk = hexutil.Encode(derivedKey.Key)
		pubkey = hexutil.Encode(derivedKey.PublicKey().Key)
	}

	// output to file or stdout
	outputContent = fmt.Sprintf("secret key: %s\npubkey: %s\n", sk, pubkey)
	if withSeed {
		outputContent += fmt.Sprintf("seed: %s\n", hexutil.Encode(seed))
	}

	if output != "" {
		if _, err := os.Stat(output); errors.Is(err, os.ErrNotExist) {
			os.WriteFile(output, []byte(outputContent), 0600)
		} else {
			fmt.Fprintf(os.Stderr, "file %s already exists, quit", output)
		}
	} else {
		fmt.Print(outputContent)
	}
}

func _derive(key *bip32.Key, path string) (*bip32.Key, error) {
	if !strings.HasPrefix(path, "m/") {
		return nil, errors.New("invalid path prefix")
	}

	path = strings.TrimPrefix(path, "m/")
	pathIndice := strings.Split(path, "/")
	childKey := key
	for _, pathIndex := range pathIndice {
		quote := false
		if strings.HasSuffix(pathIndex, "'") {
			pathIndex = strings.TrimRight(pathIndex, "'")
			quote = true
		}
		index, err := strconv.Atoi(pathIndex)
		if err != nil {
			return nil, errors.New("invalid path field")
		}

		if quote {
			index += QUOTE_PREFIX
		}
		childKey, err = childKey.NewChildKey(uint32(index))
		if err != nil {
			return nil, err
		}
	}

	return childKey, nil
}
