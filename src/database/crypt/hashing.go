package crypt

import (
	"golang.org/x/crypto/bcrypt"
)

func Hash(in string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(in), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CompareHashToPlain(hash, real string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(real)) == nil
}
