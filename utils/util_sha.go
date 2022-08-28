package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

type CryptoSha256 struct {
}

func (c CryptoSha256) Encrypt(plainText string) string {
	data := []byte(plainText)
	hash := sha256.New()
	hash.Write(data)
	v := hash.Sum(nil)
	vstr := hex.EncodeToString(v)
	return vstr
}
