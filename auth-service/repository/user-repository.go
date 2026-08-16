package repository

import (
	"context"
	"database/sql"
	"emoji-survivors/auth-service/models"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

var ErrUsernameTaken = errors.New("имя пользователя уже занято")
var ErrUserNotFound = errors.New("пользователь не найден")

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) Create(ctx context.Context, username, passwordHash string) (int, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("не удалось начать транзакцию: %w", err)
	}
	defer tx.Rollback()

	var userID int
	err = tx.QueryRowContext(ctx, `INSERT INTO users (username, password_hash)
		VALUES ($1, $2) RETURNING id`, username, passwordHash).Scan(&userID)
	if err != nil {
		if isUniqueViolation(err) {
			return 0, ErrUsernameTaken
		}
		return 0, fmt.Errorf("не удалось создать пользователя: %w", err)
	}

	_, err = tx.ExecContext(ctx, `INSERT INTO scores (user_id)
		VALUES ($1)`, userID)
	if err != nil {
		return 0, fmt.Errorf("не удалось создать запись очков: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("не удалось зафиксировать транзакцию: %w", err)
	}

	return userID, nil
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	var u models.User
	err := r.db.QueryRowContext(ctx,
		`SELECT id, username, password_hash, created_at FROM users WHERE username = $1`,
		username,
	).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("ошибка поиска пользователя: %w", err)
	}

	return &u, nil
}

func isUniqueViolation(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr.Code == "23505"
	}
	return false
}
