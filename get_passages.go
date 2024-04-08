package main

import (
	"context"
	"net/http"
	"time"

	"main/domain_model"
	"main/view"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func LoadPassagesPageModel(conn *pgxpool.Pool, user_id int32, clientDate time.Time, page interface{},) (view.AppModel, error) {
	type passageModel struct {
		Id           int
		Book         string
		StartChapter int
		StartVerse   int
		EndChapter   int
		EndVerse     int
		ReviewAt     *time.Time
	}

    query := `SELECT first_name, last_name FROM "user" WHERE id = $1`
    rows, _ := conn.Query(context.Background(), query, user_id)
    user, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[view.UserModel])
    if err != nil {
        return view.AppModel{},err
    }

    query = `
        SELECT id, book, start_chapter, start_verse, end_chapter, end_verse, review_at
        FROM passage
        WHERE user_id = $1
        ORDER BY book, start_chapter, start_verse, end_chapter, end_verse
    `
    rows, _ = conn.Query(context.Background(), query, user_id)
    dbpassages, err := pgx.CollectRows(rows, pgx.RowToStructByName[passageModel])
    if err != nil {
        return view.AppModel{},err
    }

    passages := make([]view.PassageListItemModel, len(dbpassages))
    for i, dbpassage := range dbpassages {
        passageData := view.PassageListItemModel{
            Id:        dbpassage.Id,
            Reference: domain_model.PassageReference{dbpassage.Book, dbpassage.StartChapter, dbpassage.StartVerse, dbpassage.EndChapter, dbpassage.EndVerse}.String(),
        }
        if dbpassage.ReviewAt != nil {
            passageData.ReviewAt = dbpassage.ReviewAt.Format("01-02-2006")
            passageData.ReviewDue = clientDate.Compare(*dbpassage.ReviewAt) >= 0
        }
        passages[i] = passageData
    }

    return view.AppModel {
        Page: view.PassagesPageModel{
            Passages: passages,
            Page: page,
        },
        User: &user,
    }, nil
}

func GetPassages(router *mux.Router, ctx *ServerContext) {
	type passageModel struct {
		Id           int
		Book         string
		StartChapter int
		StartVerse   int
		EndChapter   int
		EndVerse     int
		ReviewAt     *time.Time
	}

	router.HandleFunc("/passages", func(w http.ResponseWriter, r *http.Request) {
		session, err := GetSession(r, ctx)
		if err != nil {
			http.Error(w, "Session Error", http.StatusInternalServerError)
			return
		}
		if session == nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

        view.CreateViewEngine(ctx.Conn, r.Context(), w).RenderPassages(int(*session.user_id), GetClientDate(r))
	}).Methods("GET")
}
