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
	deriveCmd.Flags().Uint("start", 0, "start index when hd path contains place holder(x)")
	deriveCmd.Flags().Uint("count", 1, "derive count when hd path contains place holder(x)")

	deriveCmd.Flags().StringP("output", "o", "", "output file, print to stdout if empty")
	deriveCmd.Flags().Bool("json", false, "output json format")

}

type Key struct {
	Path      string `json:"path"`
	SecretKey string `json:"secret_key"`
	PublicKey string `json:"public_key"`
}

type OutputKey struct {
	Seed string `json:"seed,omitempty"`
	Keys []Key  `json:"keys"`
}

func (outputKey *OutputKey) JsonString() (string, error) {
	bytes, err := json.Marshal(outputKey)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (outputKey *OutputKey) String() string {
	var result string
	if outputKey.Seed != "" {
		result += fmt.Sprintf("seed: %s\n", outputKey.Seed)
	}
	for _, key := range outputKey.Keys {
		result += fmt.Sprintf("path: %s\nsecret key: %s\npublic key: %s\n", key.Path, key.SecretKey, key.PublicKey)
	}
	return result
}

func derive(cmd *cobra.Command, args []string) {
	var (
		outputContent string
		password      string
		err           error
	)

	output, _ := cmd.Flags().GetString("output")
	mnemonic, _ := cmd.Flags().GetString("mnemonic")
	mnemonic, err = utils.ReadSecret(mnemonic, "Enter mnemonic: ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "read mnemonic error: %s", err)
		os.Exit(1)
	}

	jsonFormat, _ := cmd.Flags().GetBool("json")
	enterPass, _ := cmd.Flags().GetBool("enter-pass")
	if enterPass {
		password, err = utils.ReadSecret(password, "Enter password: ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "read password error: %s", err)
			os.Exit(1)
		}
	} else {
		password, _ = cmd.Flags().GetString("password")
	}

	path, _ := cmd.Flags().GetString("path")
	start, _ := cmd.Flags().GetUint("start")
	count, _ := cmd.Flags().GetUint("count")
	seed := bip39.NewSeed(mnemonic, password)

	rootKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Create master key error: %s\n", err)
		os.Exit(1)
	}

	outputKey := OutputKey{Seed: hexutil.Encode(seed)}

	// return root sk when path is empty
	if path == "" {
		sk := hexutil.Encode(rootKey.Key)
		pubkey := hexutil.Encode(rootKey.PublicKey().Key)
		outputKey.Keys = append(outputKey.Keys, Key{Path: path, SecretKey: sk, PublicKey: pubkey})
	} else {
		keys, paths, err := utils.DerivesByPath(rootKey, path, start, count)
		if err != nil {
			fmt.Fprintf(os.Stderr, "derive error: %s\n", err)
			os.Exit(1)
		}
		for i := range keys {
			sk := hexutil.Encode(keys[i].Key)
			pubkey := hexutil.Encode(keys[i].PublicKey().Key)
			outputKey.Keys = append(outputKey.Keys, Key{Path: paths[i], SecretKey: sk, PublicKey: pubkey})
		}
	}

	if jsonFormat {
		outputContent, err = outputKey.JsonString()
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}
	} else {
		outputContent = outputKey.String()
	}

	if output != "" {
		err = utils.WriteFileWhenNotExists(output, []byte(outputContent), 0600)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}
	} else {
		fmt.Print(outputContent)
	}
}
