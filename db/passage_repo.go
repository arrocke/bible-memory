package db

import (
	"context"
	"main/domain_model"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PassageRepo interface {
	Get(id uint) (domain_model.Passage, error)
	Commit(*domain_model.Passage) error
}

type PgPassageRepo struct {
	pool *pgxpool.Pool
}

func CreatePgPassageRepo(pool *pgxpool.Pool) PgPassageRepo {
	return PgPassageRepo{pool}
}

func (repo PgPassageRepo) Get(id uint) (domain_model.Passage, error) {
	type passageModel struct {
		Id           uint
		Book         string
		StartChapter uint
		StartVerse   uint
		EndChapter   uint
		EndVerse     uint
		Text         string
		Owner        uint
		Interval     *uint
		ReviewedAt   *time.Time
		NextReview   *time.Time
	}

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

	passage := domain_model.Passage{
		Id: dbModel.Id,
		Reference: domain_model.PassageReference{
			Book:         dbModel.Book,
			StartChapter: dbModel.StartChapter,
			StartVerse:   dbModel.StartVerse,
			EndChapter:   dbModel.EndChapter,
			EndVerse:     dbModel.EndVerse,
		},
		Text:  dbModel.Text,
		Owner: dbModel.Owner,
	}
	if dbModel.Interval != nil && dbModel.ReviewedAt != nil && dbModel.NextReview != nil {
		passage.Review = &domain_model.PassageReview{
            Interval:   *dbModel.Interval,
			ReviewedAt: dbModel.ReviewedAt,
            NextReview: *dbModel.NextReview,
		}
	}

	return passage, nil
}

func (repo PgPassageRepo) Commit(passage *domain_model.Passage) error {
    query := ""
    if passage.IsNew() {
        query = `
            INSERT INTO passage (id, book, start_chapter, start_verse, end_chapter, end_verse, text, user_id, interval, reviewed_at, review_at)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
        `
    } else {
        query = `
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
    }

    var interval *uint
    var reviewedAt *time.Time
    var nextReview *time.Time
    if passage.Review != nil {
        interval = &passage.Review.Interval
        reviewedAt = passage.Review.ReviewedAt
        nextReview = &passage.Review.NextReview
    }

    _, err := repo.pool.Exec(
        context.Background(), query,
        passage.Id,
        passage.Reference.Book,
        passage.Reference.StartChapter,
        passage.Reference.StartVerse,
        passage.Reference.EndChapter,
        passage.Reference.EndVerse,
        passage.Text,
        passage.Owner,
        interval,
        reviewedAt,
        nextReview,
    )

	return err
}
