package utils

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateSalt(length int) (string, error) {
	buf := make([]byte, length)

	if _, err := rand.Read(buf); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(buf), nil
}
