package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

func GenerateStringSalt(length int) (string, error) {
	bytes := make([]byte, length)

	n, err := rand.Read(bytes)

	if n != length || err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

func HashString(s string, salt string) string {
	hasher := sha256.New()

	hasher.Write([]byte(s + salt))

	return hex.EncodeToString(hasher.Sum(nil))
}

func VerifyString(s string, salt string, hash string) bool {
	newHash := HashString(s, salt)
	return newHash == hash
}
