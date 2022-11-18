package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/tyler-smith/go-bip32"
)

const QUOTE_PREFIX = 0x80000000

func DeriveByPath(key *bip32.Key, path string) (*bip32.Key, error) {
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

// path format: m/60'/44'/x/0
// x is place holder to derive from start to start + count - 1  secret key
func DerivesByPath(key *bip32.Key, path string, start, count uint) (keys []*bip32.Key, paths []string, err error) {
	if strings.Count(path, "x") == 0 {
		key, err = DeriveByPath(key, path)
		keys = append(keys, key)
		paths = append(paths, path)
		return
	}

	if strings.Count(path, "x") != 1 {
		return nil, nil, errors.New("invalid path format,path must has one x")
	}

	newPath := strings.Replace(path, "x", "%d", 1)
	for i := start; i < start+count; i++ {
		path = fmt.Sprintf(newPath, i)
		key, err := DeriveByPath(key, path)
		if err != nil {
			return nil, nil, err
		}
		keys = append(keys, key)
		paths = append(paths, path)
	}
	return
}

func t() {
}
