package commonUtils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// StringSliceToString converts a slice of strings to a comma delimited string
func StringSliceToString(src []string) (tgt string) {
	tgt = "["
	for i := 0; i < len(src); i++ {
		joinChr := ", "
		if i == 0 {
			joinChr = ""
		}
		tgt = fmt.Sprintf("%s%s\"%s\"", tgt, joinChr, src[i])
	}
	tgt = fmt.Sprintf("%s]", tgt)
	return
}

// GenerateRandomString returns a securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}
