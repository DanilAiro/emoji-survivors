package repository

import (
	"database/sql"
	"emoji-survivors/auth-service/db"
)

var database *sql.DB

func ConnectDB() {
	database = db.ConnectDB()
}

func AddNewUserToUsers(login, password string) (int, error) {
	query := `
		INSERT INTO users (login, password_hash)
		VALUES ($1, $2)
	`

	res, err := database.Exec(query, login, password)
	if err != nil {
		return 0, err
	}
	
	user_id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(user_id), nil
}

func AddNewUserToScores(user_id int) error {
	query := `
		INSERT INTO scores (user_id)
		VALUES ($1)
	`

	_, err := database.Exec(query, user_id)

	return err
}

func AddNewScoreIfHigher(user_id, score int) error {
	query := `
		UPDATE INTO scores (user_id, score, updated_at)
		VALUES ($1, $2, NOW())
	`

	_, err := database.Exec(query, user_id, score)

	return err
}