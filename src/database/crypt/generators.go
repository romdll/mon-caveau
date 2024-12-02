package crypt

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"moncaveau/utils"
	"strconv"
	"strings"
)

const AccountKeyLength = 18

func GenerateSecureAccountKey() (string, error) {
	logger.Println("Generating secure account key...")

	upperLimit := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(AccountKeyLength)), nil)

	randomNum, err := rand.Int(rand.Reader, upperLimit)
	if err != nil {
		logger.Printf("Error generating random number: %v", err)
		return "", err
	}

	accountKey := fmt.Sprintf("%0"+strconv.Itoa(AccountKeyLength)+"d", randomNum)

	var result []string
	for i := 0; i < len(accountKey); i += 6 {
		end := i + 6
		if end > len(accountKey) {
			end = len(accountKey)
		}
		result = append(result, accountKey[i:end])
	}
	accountKeyWithHyphens := strings.Join(result, "-")

	logger.Printf("Generated account key: %s", utils.MaskOnlyNumbers(accountKeyWithHyphens, 6))

	return accountKeyWithHyphens, nil
}
