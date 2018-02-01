package utils

import (
	"crypto/sha256"
)

func CalcSha256Hash(source []byte) []byte {
	hash := sha256.New()
	hash.Write(source)
	md := hash.Sum(nil)
	return md
}
