package repository

import (
	"context"
	"database/sql"
	"fmt"
)

type ScoreRepository struct {
	db *sql.DB
}

func NewScoreRepository(db *sql.DB) *ScoreRepository {
	return &ScoreRepository{db}
}

func (r *ScoreRepository) UpdateIfHigher(ctx context.Context, user_id, score int) error {
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE scores SET score = $2, updated_at = now()
     	WHERE user_id = $1 AND score < $2`,
		user_id,
		score)
	if err != nil {
		return fmt.Errorf("не удалось обновить счёт игрока %d: %w", user_id, err)
	}

	return err
}
