package pg

import (
	"database/sql"
	"encoding/hex" // Needed for encoding
	"fmt"
)

var DB *sql.DB

func InsertInDB(username, email string, passwordHash, salt []byte) error {
	// Encode data to hex strings
	hashHex := hex.EncodeToString(passwordHash)
	saltHex := hex.EncodeToString(salt)

	// Insert the hex strings as plain text
	_, err := DB.Exec(
		"INSERT INTO users (username, email, password_hash, salt) VALUES ($1, $2, $3, $4)",
		username,
		email,
		hashHex,
		saltHex,
	)
	if err != nil {
		return fmt.Errorf("failed to insert user into database: %w", err)
	}

	return nil
}

