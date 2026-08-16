package db

import (
	"database/sql"
	"log"
)

func ConnectDB(DBUrl string) *sql.DB {
	db, err := sql.Open("postgres", DBUrl)

	if (err != nil) {
		log.Fatal("неудалось подключиться к базе данных")
	}

	err = db.Ping()

	if (err != nil) {
		log.Fatal("база данных недоступна")
	}

	return db
}
