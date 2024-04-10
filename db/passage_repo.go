package db

import (
	"context"
	"main/domain_model"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PassageRepo interface {
	Get(id int) (domain_model.Passage, error)
	Create(*domain_model.Passage) error
	Update(*domain_model.Passage) error
}

type PgPassageRepo struct {
	pool *pgxpool.Pool
}

func CreatePgPassageRepo(pool *pgxpool.Pool) PgPassageRepo {
	return PgPassageRepo{pool}
}

type passageModel struct {
	Id           int
	Book         string
	StartChapter int
	StartVerse   int
	EndChapter   int
	EndVerse     int
	Text         string
	Owner        int
	Interval     *int
	ReviewedAt   *time.Time
	NextReview   *time.Time
}

func (dbModel *passageModel) toDomain() domain_model.Passage {
	var reviewState *domain_model.PassageReview
	if dbModel.Interval != nil && dbModel.NextReview != nil {
		reviewState = &domain_model.PassageReview{
			Interval:   domain_model.ReviewInterval(*dbModel.Interval),
			NextReview: domain_model.ReviewTimestamp(*dbModel.NextReview),
			ReviewedAt: (*domain_model.ReviewTimestamp)(dbModel.NextReview),
		}
	}
	passage := domain_model.Passage{
		Id: dbModel.Id,
		Reference: domain_model.PassageReference{
			Book:         dbModel.Book,
			StartChapter: dbModel.StartChapter,
			StartVerse:   dbModel.StartVerse,
			EndChapter:   dbModel.EndChapter,
			EndVerse:     dbModel.EndVerse,
		},
		Text:        dbModel.Text,
		Owner:       dbModel.Owner,
		ReviewState: reviewState,
	}

	return passage
}

func passageToDb(passage *domain_model.Passage) passageModel {
	dbModel := passageModel{
		Id:           passage.Id,
		Book:         passage.Reference.Book,
		StartChapter: passage.Reference.StartChapter,
		StartVerse:   passage.Reference.StartVerse,
		EndChapter:   passage.Reference.EndChapter,
		EndVerse:     passage.Reference.EndVerse,
		Text:         passage.Text,
		Owner:        passage.Owner,
	}

	if passage.ReviewState != nil {
		dbModel.Interval = (*int)(&passage.ReviewState.Interval)
		dbModel.ReviewedAt = (*time.Time)(passage.ReviewState.ReviewedAt)
		dbModel.NextReview = (*time.Time)(&passage.ReviewState.NextReview)
	}

	return dbModel
}

func (repo PgPassageRepo) Get(id int) (domain_model.Passage, error) {
	query := `
        SELECT id, book, start_chapter, start_verse, end_chapter, end_verse, text, user_id, interval, reviewed_at, review_at
        FROM passage
        WHERE id = $1
    `
	rows, _ := repo.pool.Query(context.Background(), query, id)
	dbModel, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByPos[passageModel])
	if err != nil {
		return domain_model.Passage{}, err
	}

	passage := dbModel.toDomain()
	return passage, nil
}

func (repo PgPassageRepo) Create(passage *domain_model.Passage) error {
	query := `
        INSERT INTO passage (id, book, start_chapter, start_verse, end_chapter, end_verse, text, user_id, interval, reviewed_at, review_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
    `
	dbModel := passageToDb(passage)
	_, err := repo.pool.Exec(
		context.Background(),
        query,
		dbModel.Id,
		dbModel.Book,
		dbModel.StartChapter,
		dbModel.StartVerse,
		dbModel.EndChapter,
		dbModel.EndVerse,
		dbModel.Text,
		dbModel.Owner,
		dbModel.Interval,
		dbModel.ReviewedAt,
		dbModel.NextReview,
	)
	return err
}

func (repo PgPassageRepo) Update(passage *domain_model.Passage) error {
	query := `
        UPDATE passage
            SET book = $2,
            start_chapter = $3,
            start_verse = $4,
            end_chapter = $5,
            end_verse = $6,
            text = $7,
            user_id = $8,
            interval = $9,
            reviewed_at = $10,
            review_at = $11
        WHERE id = $1
    `
	dbModel := passageToDb(passage)
	_, err := repo.pool.Exec(
		context.Background(),
        query,
		dbModel.Id,
		dbModel.Book,
		dbModel.StartChapter,
		dbModel.StartVerse,
		dbModel.EndChapter,
		dbModel.EndVerse,
		dbModel.Text,
		dbModel.Owner,
		dbModel.Interval,
		dbModel.ReviewedAt,
		dbModel.NextReview,
	)
	return err
}
