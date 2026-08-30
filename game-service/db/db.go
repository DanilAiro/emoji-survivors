package db

import (
	"database/sql"
	"errors"

	_ "github.com/lib/pq"
)

var ErrNoConnect = errors.New("неудалось подключиться к базе данных")
var ErrNoPing = errors.New("база данных недоступна")

func ConnectDB(DBUrl string) (*sql.DB, error) {
	db, err := sql.Open("postgres", DBUrl)

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
