package repository

import (
	"context"
	"database/sql"
	"emoji-survivors/auth-service/models"
	"fmt"
)

type ScoreRepository struct {
	db *sql.DB
}

func NewScoreRepository(db *sql.DB) *ScoreRepository {
	return &ScoreRepository{db}
}

func (r *ScoreRepository) GetTop10(ctx context.Context) ([]models.ScoreWithUsername, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT u.username, s.score FROM scores s JOIN users u ON u.id = s.user_id ORDER BY s.score DESC LIMIT 10`)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить топ-10: %w", err)
	}
	defer rows.Close()

	var result []models.ScoreWithUsername
	for rows.Next() {
		var s models.ScoreWithUsername
		if err := rows.Scan(&s.Username, s.Score); err != nil {
			return nil, fmt.Errorf("не удалось прочитать строку результата: %w", err)
		}
		result = append(result, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при переборе строк: %w", err)
	}

	return result, nil
}
