package core

import (
	"crypto/sha256"
	"encoding/hex"
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

func Shorten(inp string) string {
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
				hoop := int(url[i] - '0') // Convert byte to integer directly
				if hoop == 0 {            // Prevent division by zero
					hoop = 1
				}
				i += len(url) % hoop
				if i >= len(url) { // Avoid out-of-bounds access
					break
				}
				continue
			}
			i += 2
		}
		url = res
	}
	return hex.EncodeToString(url)
}
