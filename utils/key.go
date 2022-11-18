package utils

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"errors"
	"fmt"
	"math/big"

	btcutil "github.com/FactomProject/btcutilecc"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	curve = btcutil.Secp256k1()
	// curveParams = curve.Params()
)

const (
	PublicKeyCompressedLength = 33
)

func PublieKey(privatekey string) (string, error) {
	priKey, err := hexutil.Decode(privatekey)
	if err != nil {
		return "", err
	}
	pubkey := PublicKeyForPrivateKey(priKey)
	return hexutil.Encode(pubkey), nil
}

func PublicKeyForPrivateKey(key []byte) []byte {
	return CompressPublicKey(curve.ScalarBaseMult(key))
}

func CompressPublicKey(x *big.Int, y *big.Int) []byte {
	var key bytes.Buffer

	// Write header; 0x2 for even y value; 0x3 for odd
	key.WriteByte(byte(0x2) + byte(y.Bit(0)))

	// Write X coord; Pad the key so x is aligned with the LSB. Pad size is key length - header size (1) - xBytes size
	xBytes := x.Bytes()
	for i := 0; i < (PublicKeyCompressedLength - 1 - len(xBytes)); i++ {
		key.WriteByte(0x0)
	}
	key.Write(xBytes)

	return key.Bytes()
}

func IsCompressedPublicKey(pubkey []byte) (bool, error) {
	if len(pubkey) == 33 {
		if pubkey[0] == 2 || pubkey[0] == 3 {
			return true, nil
		} else {
			return false, errors.New("invalid pubkey")
		}
	} else if len(pubkey) == 65 {
		if pubkey[0] == 4 {
			return false, nil
		} else {
			return false, errors.New("invalid pubkey")
		}
	}
	return false, errors.New("invalid pubkey")
}

// Not Working
func RecoverPublicKeyFromCompressed(pubkey []byte) (*ecdsa.PublicKey, error) {
	compressed, err := IsCompressedPublicKey(pubkey)
	if err != nil {
		return nil, err
	}
	if !compressed {
		return nil, errors.New("not compressed public key")
	}

	fmt.Printf("pubkey: %+v\n", pubkey)
	publicKey := new(ecdsa.PublicKey)
	x, y := elliptic.UnmarshalCompressed(crypto.S256(), pubkey)
	if x == nil || y == nil {
		return nil, errors.New("unmarshal compressed public key error")
	}
	fmt.Printf("pubkey: %+v\n", pubkey)
	fmt.Printf("x: %+v\n", x)
	fmt.Printf("y: %+v\n", y)

	publicKey.Curve = crypto.S256()
	publicKey.X = x
	publicKey.Y = y
	return publicKey, nil
}
