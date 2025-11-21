package pg

import (
	"database/sql"
	"fmt"
)

var DB *sql.DB

// InsertInDB inserts a new user into the database.
// It now includes email and returns an error on failure.
func InsertInDB(username, email string, passwordHash, salt []byte) error {
	_, err := DB.Exec(
		"INSERT INTO users (username, email, password_hash, salt) VALUES ($1, $2, $3, $4)",
		username,
		email,
		passwordHash,
		salt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert user into database: %w", err)
	}

	return nil
}

