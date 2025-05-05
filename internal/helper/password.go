package helper

import (
	"crypto/rand"
	"math/big"
)

func PasswordGenerate(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)

	m := big.NewInt(int64(len(charset)))

	for i := range result {
		for {
			n, err := rand.Int(rand.Reader, m)
			if err != nil {
				continue
			}

			result[i] = charset[n.Int64()]

			break
		}
	}

	return string(result)
}
