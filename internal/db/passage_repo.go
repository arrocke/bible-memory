package db

import (
	"context"
	"main/internal/model"

	"github.com/georgysavva/scany/v2/pgxscan"
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
    var passage model.Passage
    err := pgxscan.Select(c, r.Pool, &passage, query, id)
	return passage, err
}

func (r PassageRepo) GetPassagesForOwner(c context.Context, ownerId int) ([]model.Passage, error) {
    query := `
        SELECT id, book, start_chapter, start_verse, end_chapter, end_verse, text, user_id, interval, reviewed_at, review_at
        FROM passage
        WHERE user_id = $1
    `
    var passages []model.Passage
    err := pgxscan.Select(c, r.Pool, &passages, query, ownerId)
	return passages, err
}

func (r PassageRepo) Create(c context.Context, passage model.Passage) error {
    query := `
        INSERT INTO passage (book, start_chapter, start_verse, end_chapter, end_verse, text, user_id, interval, reviewed_at, review_at)
        VALUES (@book, @start_chapter, @start_verse, @end_chapter, @end_verse, @text, @user_id, @interval, @reviewed_at, @review_at)
    `
    _, err := r.Pool.Exec(c, query, pgx.NamedArgs{
        "book": passage.Reference.Book,
        "start_chapter": passage.Reference.StartChapter,
        "start_verse": passage.Reference.StartVerse,
        "end_chapter": passage.Reference.EndChapter,
        "end_verse": passage.Reference.EndVerse,
        "text": passage.Text,
        "user_id": passage.Owner,
        "interval": passage.Interval,
        "reviewed_at": passage.ReviewedAt,
        "review_at": passage.NextReview,
    })
	return err
}

func (r PassageRepo) Update(c context.Context, passage model.Passage) error {
    query := `
        UPDATE passage SET 
            book = @book,
            start_chapter = @start_chapter,
            start_verse = @start_verse,
            end_chapter = @end_chapter,
            end_verse = @end_verse,
            text = @text,
            user_id = @user_id,
            interval = @interval,
            reviewed_at = @reviewed_at,
            review_at = @review_at
        WHERE id = @id
    `
    _, err := r.Pool.Exec(c, query, pgx.NamedArgs{
        "id": passage.Id,
        "book": passage.Reference.Book,
        "start_chatper": passage.Reference.StartChapter,
        "start_verse": passage.Reference.StartVerse,
        "end_chapter": passage.Reference.EndChapter,
        "end_verse": passage.Reference.EndVerse,
        "text": passage.Text,
        "user_id": passage.Owner,
        "interval": passage.Interval,
        "reviewed_at": passage.ReviewedAt,
        "review_at": passage.NextReview,
    })
	return err
}
