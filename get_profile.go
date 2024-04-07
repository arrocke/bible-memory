package main

import (
	"main/view"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

func GetProfile(router *mux.Router, ctx *ServerContext) {
	type userModel struct {
		Email     string
		FirstName string
		LastName  string
	}

	router.HandleFunc("/users/profile", func(w http.ResponseWriter, r *http.Request) {
		session, err := GetSession(r, ctx)
		if err != nil {
			http.Error(w, "Session Error", http.StatusInternalServerError)
			return
		}
		if session == nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

        query := `SELECT email, first_name, last_name FROM "user" WHERE id = $1`
        rows, _ := ctx.Conn.Query(r.Context(), query, session.user_id)
        user, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[userModel])
		if err != nil {
			http.Error(w, "Database Error", http.StatusInternalServerError)
			return
		}

        view.App(view.AppModel {
            Page: view.ProfilePageModel{
                Email: user.Email,
                FirstName: user.FirstName,
                LastName: user.LastName,
            },
            User: &view.UserModel{
                FirstName: user.FirstName,
                LastName: user.LastName,
            },
        }).Render(r.Context(), w)
	}).Methods("Get")
}
