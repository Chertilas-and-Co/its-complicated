package pg

import (
  "log"
  "database/sql"
)

var db *sql.DB

func insertInDB(username string, passwordHash []byte, salt []byte) {
	result, err := db.Exec("INSERT INTO users (username, password_hash, salt) VALUES ($1, $2, $3)", username, passwordHash[:], salt)
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

