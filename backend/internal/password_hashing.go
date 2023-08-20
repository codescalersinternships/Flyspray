package internal

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"io"
)

const saltLen = 16

// HashPassword hashes password of user
func HashPassword(password []byte) ([]byte, error) {
	salt := make([]byte, saltLen)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return []byte{}, err
	}

	hashedPassword := sha256.Sum256(append(salt, password...))
	return append(salt, hashedPassword[:]...), nil
}

// CheckPasswordHash checks if given password is same as hashed one
func CheckPasswordHash(hashedPassword []byte, password string) bool {
	hashedPasswordCopy := make([]byte, len(hashedPassword))

	copy(hashedPasswordCopy, hashedPassword)
	salt := hashedPasswordCopy[:saltLen]

	checkedPass := sha256.Sum256(append(salt, []byte(password)...))
	return bytes.Equal(append(salt, checkedPass[:]...), hashedPassword)
}
