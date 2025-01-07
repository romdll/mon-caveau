package utils

import (
	"crypto/rand"
	"math/big"
)

func SafeIntN(max int) int {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return -1
	}
	return int(n.Int64())
}
