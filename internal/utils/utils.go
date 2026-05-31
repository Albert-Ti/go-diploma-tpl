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
	numStr := strings.Split(order, "")
	sum := 0

	for i, v := range numStr {
		n, err := strconv.Atoi(v)
		if err != nil {
			panic(err)
		}
		if i%2 == 0 {
			s := n * 2
			if s > 10 {
				res := 0
				slice := strings.SplitSeq(strconv.Itoa(s), "")
				for v := range slice {
					n2, _ := strconv.Atoi(v)
					res += n2
				}
				sum += res
			} else {
				sum += s
			}
		} else {
			sum += n
		}
	}
	return sum%10 == 0
}
