package services

import "emoji-survivors/auth-service/repository"

func CreateNewUser(login, password_hash string) error {
	user_id, err := repository.AddNewUserToUsers(login, password_hash)
	if err != nil {
		return err
	}

	err = repository.AddNewUserToScores(user_id)

	return err
}