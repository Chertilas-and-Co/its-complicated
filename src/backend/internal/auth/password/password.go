package password

import (
	"crypto/rand"
	"crypto/sha256"
)

func generateSalt(length int) ([]byte, error) {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

func hashPassword(password string, salt []byte) (hash []byte) {
	data := append([]byte(password), salt...)
	hash1 := sha256.Sum256(data)
	hash = hash1[:]
	return hash
}
