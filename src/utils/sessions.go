package utils

import (
	"crypto/rand"
	"encoding/base64"
)

const (
	sessionTokenBytesSize = 64
)

func GenerateSessionToken() (string, error) {
	tokenBytes := make([]byte, sessionTokenBytesSize)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(tokenBytes), nil
}
