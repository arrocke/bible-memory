package main

import (
	"context"
	"main/view"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func LoadPassages(context context.Context, conn *pgxpool.Pool, userId int, clientDate time.Time) ([]view.PassageListItemModel, error) {
	type passageModel struct {
		Id           int
		Book         string
		StartChapter int
		StartVerse   int
		EndChapter   int
		EndVerse     int
		ReviewAt     *time.Time
	}

	query := `
        SELECT id, book, start_chapter, start_verse, end_chapter, end_verse, review_at
        FROM passage
        WHERE user_id = $1
        ORDER BY book, start_chapter, start_verse, end_chapter, end_verse
    `
	rows, _ := conn.Query(context, query, userId)
	dbpassages, err := pgx.CollectRows(rows, pgx.RowToStructByName[passageModel])
	if err != nil {
		return nil, err
	}

	passages := make([]view.PassageListItemModel, len(dbpassages))
	for i, dbpassage := range dbpassages {
		passageData := view.PassageListItemModel{
			Id: dbpassage.Id,
			Reference: view.PassageReference{
				Book:         dbpassage.Book,
				StartChapter: dbpassage.StartChapter,
				StartVerse:   dbpassage.StartVerse,
				EndChapter:   dbpassage.EndChapter,
				EndVerse:     dbpassage.EndVerse,
			}.String(),
		}
		if dbpassage.ReviewAt != nil {
			passageData.ReviewAt = dbpassage.ReviewAt.Format("01-02-2006")
			passageData.ReviewDue = clientDate.Compare(*dbpassage.ReviewAt) >= 0
		}
		passages[i] = passageData
	}

	return passages, nil
}

func (ctx *ServerContext) passagesRoutes(router *http.ServeMux) {
	router.Handle("GET /passages", AuthMiddleware(true, HandleErrors(func(w http.ResponseWriter, r *http.Request) error {
		userId := GetUserId(r)
		clientDate := GetClientDate(r)

		passages, err := LoadPassages(r.Context(), ctx.Conn, userId, clientDate)
		if err != nil {
			return err
		}

		page := view.PassagesPage(view.PassagesPageModel{Passages: passages})
        return ctx.RenderPage(w, r, page)
	})))
}
