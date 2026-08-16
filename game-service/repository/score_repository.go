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

// обновить до UpdateIfHigher
func (r *ScoreRepository) Update(ctx context.Context, user_id, score int) error {
	_, err := r.db.ExecContext(ctx, `UPDATE INTO scores (user_id, score, updated_at)
		VALUES ($1, $2, NOW())`, user_id, score)
	if err != nil {
		return fmt.Errorf("не удалось обновить счёт игрока %d: %w", user_id, err)
	}

	return err
}
