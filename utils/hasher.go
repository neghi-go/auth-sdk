package utils

import (
	"crypto/subtle"
	"encoding/base64"
	"errors"

	"golang.org/x/crypto/argon2"
)

var (
	ErrMismatch = errors.New("")
)

type Hasher struct {
	Time    uint32
	Memory  uint32
	Threads uint8
	KeyLen  uint32
}

func NewHasher() *Hasher {
	return &Hasher{
		Time:    1,
		Memory:  16 * 1024,
		Threads: 4,
		KeyLen:  32,
	}
}

func (h *Hasher) Hash(password, salt string) string {
	key := argon2.IDKey([]byte(password), []byte(salt), h.Time, h.Memory, h.Threads, h.KeyLen)
	return base64.RawURLEncoding.EncodeToString(key)
}

func (h *Hasher) Compare(compare string, password, salt string) error {
	hashedPass := h.Hash(password, salt)

	if subtle.ConstantTimeCompare([]byte(compare), []byte(hashedPass)) != 1 {
		return ErrMismatch
	}
	return nil
}
