package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func ConnectDB(DBUrl string) *sql.DB {
	db, err := sql.Open("postgres", DBUrl)

	if err != nil {
		log.Fatalf("не удалось открыть подключение к базе данных: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("база данных недоступна: %v", err)
	}

	return db
}
