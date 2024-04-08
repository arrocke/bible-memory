package main

import (
	"context"
	"main/view"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TemplateUser struct {
	ID        int32
	FirstName string
	LastName  string
	Email     string
}
type LayoutTemplateData struct {
	User *TemplateUser
}

func LoadLayoutTemplateData(conn *pgxpool.Pool, user_id *int32) (*LayoutTemplateData, error) {
	if user_id == nil {
		return &LayoutTemplateData{}, nil
	}

	query := `SELECT id, email, first_name, last_name FROM "user" WHERE id = $1`
	rows, _ := conn.Query(context.Background(), query, user_id)
	defer rows.Close()

	user, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[TemplateUser])
	if err != nil {
		println(err.Error())
		return nil, err
	}

	templateData := LayoutTemplateData{
		User: &user,
	}

	return &templateData, nil
}

func GetIndex(router *mux.Router, ctx *ServerContext) {
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		session, err := GetSession(r, ctx)
		if err != nil {
			http.Error(w, "Session Error", http.StatusInternalServerError)
			return
		}

		if session == nil {
			view.CreateViewEngine(ctx.Conn, r.Context(), w).RenderIndex()
		} else {
			http.Redirect(w, r, "/passages", http.StatusFound)
		}

	}).Methods("Get")
}
