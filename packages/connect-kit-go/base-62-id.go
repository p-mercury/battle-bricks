package connectkit

import (
	"crypto/rand"
	"errors"
	"fmt"
)

const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
const base = 62
const limit = 256 - (256 % base)

func NewBase62Id(prefix string, length int) (string, error) {
	if length <= 0 {
		return "", errors.New("length must be > 0")
	}

	total := len(prefix) + length
	out := make([]byte, total)
	copy(out, prefix)

	i := len(prefix)
	buf := make([]byte, 64)

	for i < total {
		if _, err := rand.Read(buf); err != nil {
			return "", fmt.Errorf("connectkit: failed to read random bytes: %w", err)
		}
		for _, b := range buf {
			if b < limit {
				out[i] = alphabet[b%base]
				i++
				if i == total {
					break
				}
			}
		}
	}

	return string(out), nil
}
