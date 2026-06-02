package utils

import (
	"crypto/rand"
	"encoding/hex"
	"strconv"
	"strings"
)

func RandomHash(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

func AlgoLuna(order string) bool {
	order = strings.TrimSpace(order)

	var sum int
	var isSecond bool

	for i := len(order) - 1; i >= 0; i-- {
		n, err := strconv.Atoi(string(order[i]))
		if err != nil {
			return false
		}

		if isSecond {
			n *= 2
			if n >= 10 {
				n = n/10 + n%10
			}
		}

		sum += n
		isSecond = !isSecond
	}

	return sum%10 == 0
}
