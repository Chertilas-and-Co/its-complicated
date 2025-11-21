package pg

import (
	"database/sql"
	"log"
)

var DB *sql.DB

func InsertInDB(username string, passwordHash []byte, salt []byte) {
	result, err := DB.Exec("INSERT INTO users (username, email, password_hash, salt) VALUES ('aboba', $1, $2, $3)", username, passwordHash[:], salt)
	if err != nil {
		log.Println(err)
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		log.Println(err)
	}
	log.Println("Rows affected:", rowsAffected)

	lastInsertID, err := result.LastInsertId()

	if err != nil {
		log.Println(err)
	}
	log.Println("Last inserted id:", lastInsertID)

}
