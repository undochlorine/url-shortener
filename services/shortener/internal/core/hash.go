package core

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
)

var numbers = []byte("0123456789")

func isByteANumber(b byte) bool {
	for _, n := range numbers {
		if b == n {
			return true
		}
	}
	return false
}

func hash(inp string) string {
	if len(inp) == 0 {
		return ""
	}

	hash := sha256.Sum256([]byte(inp))
	url := hash[:]

	for hex.EncodedLen(len(url)) > 8 {
		res := make([]byte, 0, hex.EncodedLen(len(url)))
		for i := 0; i < len(url); {
			res = append(res, url[i])
			if isByteANumber(url[i]) {
				hoop, _ := strconv.Atoi(string(url[i]))
				i += len(url) % hoop
				continue
			}
			i += 2
		}
		url = res
	}
	return hex.EncodeToString(url)
}

func Shorten(url string) string {
	slashes := 0
	for i, r := range url {
		if r == '/' {
			slashes++
		}
		if slashes == 2 {
			hashed := hash(url[i+1:])
			if hashed == "" {
				return ""
			}
			return "http://" + hashed
		}
	}
	return ""
}
