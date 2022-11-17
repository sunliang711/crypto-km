package utils

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"syscall"

	"golang.org/x/term"
)

func WriteFileWhenNotExists(outputfile string, data []byte, perm fs.FileMode) error {
	if _, err := os.Stat(outputfile); errors.Is(err, os.ErrNotExist) {
		return os.WriteFile(outputfile, data, perm)
	} else {
		return fmt.Errorf("file %s already exists, quit", outputfile)
	}
}

func ReadSecret(val string, prompt string) (string, error) {
	if val == "" {
		fmt.Fprintf(os.Stderr, "%s", prompt)
		secret, err := term.ReadPassword(syscall.Stdin)
		if err != nil {
			return "", err
		}
		fmt.Fprint(os.Stderr, "\n")
		return string(secret), nil
	}
	return val, nil
}
