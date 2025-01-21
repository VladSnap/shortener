package helpers

import (
	crypto "crypto/rand"
	"math/big"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) (string, error) {
	b := make([]rune, n)
	maxLetters := big.NewInt(int64(len(letters)))
	for i := range b {
		rndIndex, err := crypto.Int(crypto.Reader, maxLetters)
		if err != nil {
			return "", err
		}
		b[i] = letters[rndIndex.Int64()]
	}
	return string(b), nil
}
