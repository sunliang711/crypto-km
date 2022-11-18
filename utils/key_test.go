package utils

import (
	"testing"
)

func TestKey(t *testing.T) {
	sk := "0x801af14cf9ecaf8ca1f3498e1958b3f9866a13b95a5007ae37ddc2cdd1c9f0c5"
	pk, err := PublieKey(sk)
	if err != nil {
		t.Error(err)
	}
	t.Logf(pk)
}
