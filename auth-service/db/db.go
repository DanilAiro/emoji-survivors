package db

import (
	"database/sql"
	"log"
)

func ConnectDB() *sql.DB {
	connStr := "user=* password=* dbname=* sslmode=disable"
	db, err := sql.Open("postgres", connStr)

	if (err != nil) {
		log.Fatal("неудалось подключиться к базе данных")
	}

	err = db.Ping()

	if (err != nil) {
		log.Fatal("база данных недоступна")
	}

	return db
}
