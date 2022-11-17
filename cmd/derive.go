/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/spf13/cobra"
	"github.com/sunliang711/crypto-km/utils"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

// deriveCmd represents the derive command
var deriveCmd = &cobra.Command{
	Use:   "derive",
	Short: "derive secret key from mnemonic password and path",
	Long:  `derive secret key from mnemonic password and path`,
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
	deriveCmd.Flags().Bool("enter-pass", false, "enter password in safe way")
	deriveCmd.Flags().String("path", "", "hd path")
	deriveCmd.Flags().StringP("output", "o", "", "output file, print to stdout if empty")
	deriveCmd.Flags().Bool("with-seed", false, "show seed")
	deriveCmd.Flags().Bool("json", false, "output json format")

}

type OutputKey struct {
	SecretKey string `json:"secret_key,omitempty"`
	PublicKey string `json:"public_key,omitempty"`
	Seed      string `json:"seed,omitempty"`
}

func (outputKey *OutputKey) JsonString() (string, error) {
	bytes, err := json.Marshal(outputKey)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (outputKey *OutputKey) String() string {
	result := fmt.Sprintf("secret key: %s\npublic key: %s\n", outputKey.SecretKey, outputKey.PublicKey)
	if outputKey.Seed != "" {
		result += fmt.Sprintf("seed: %s\n", outputKey.Seed)
	}
	return result
}

func derive(cmd *cobra.Command, args []string) {
	var (
		outputContent string
		sk            string
		pubkey        string
		password      string
		err           error
	)

	output, _ := cmd.Flags().GetString("output")
	mnemonic, _ := cmd.Flags().GetString("mnemonic")
	mnemonic, err = utils.ReadSecret(mnemonic, "Enter mnemonic: ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "read mnemonic error: %s", err)
		return
	}

	jsonFormat, _ := cmd.Flags().GetBool("json")
	enterPass, _ := cmd.Flags().GetBool("enter-pass")
	if enterPass {
		password, err = utils.ReadSecret(password, "Enter password: ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "read password error: %s", err)
			return
		}
	} else {
		password, _ = cmd.Flags().GetString("password")
	}

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
		derivedKey, err := utils.DeriveByPath(rootKey, path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "derive error: %s\n", err)
			return
		}
		sk = hexutil.Encode(derivedKey.Key)
		pubkey = hexutil.Encode(derivedKey.PublicKey().Key)
	}

	// output to file or stdout
	outputKey := OutputKey{SecretKey: sk, PublicKey: pubkey}

	if withSeed {
		// outputContent += fmt.Sprintf("seed: %s\n", hexutil.Encode(seed))
		outputKey.Seed = hexutil.Encode(seed)
	}

	if jsonFormat {
		outputContent, err = outputKey.JsonString()
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			return
		}
	} else {
		outputContent = outputKey.String()
	}

	if output != "" {
		err = utils.WriteFileWhenNotExists(output, []byte(outputContent), 0600)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			return
		}
	} else {
		fmt.Print(outputContent)
	}
}
