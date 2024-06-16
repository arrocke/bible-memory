package db

import (
	"context"
	"main/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PassageRepo struct {
    Pool *pgxpool.Pool
}

func (r PassageRepo) GetPassageById(c context.Context, id int) (model.Passage, error) {
    query := `
        SELECT id, book, start_chapter, start_verse, end_chapter, end_verse, text, user_id, interval, reviewed_at, review_at
        FROM passage
        WHERE id = $1
    `
	rows, _ := r.Pool.Query(context.Background(), query, id)
	passage, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Passage])
	if err != nil {
		return model.Passage{}, err
	}

	return passage, nil
}

func (r PassageRepo) GetPassagesForOwner(c context.Context, ownerId int) ([]model.Passage, error) {
    query := `
        SELECT id, book, start_chapter, start_verse, end_chapter, end_verse, text, user_id, interval, reviewed_at, review_at
        FROM passage
        WHERE user_id = $1
    `
	rows, _ := r.Pool.Query(context.Background(), query, ownerId)
	passages, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Passage])
	if err != nil {
		return nil, err
	}

	return passages, nil
}
