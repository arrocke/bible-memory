package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

func PostLogin(router *mux.Router, ctx *ServerContext) {
	type LoginForm struct {
		Email    string
		Password string
	}

	type DbUser struct {
		ID       int32
		Password string
	}

	router.HandleFunc("/users/login", func(w http.ResponseWriter, r *http.Request) {
		form := LoginForm{
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
		}

		query := `SELECT id, password FROM "user" WHERE email = $1`
		rows, _ := ctx.Conn.Query(context.Background(), query, form.Email)
		user, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[DbUser])
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				LoginTemplate.ExecuteTemplate(w, "layout.html", LoginTemplateData{Email: form.Email, Error: "Invalid email or password"})
			} else {
				println(err.Error())
				http.Error(w, "Database Error", http.StatusInternalServerError)
			}
			return
		}

		if user.Password != form.Password {
			LoginTemplate.ExecuteTemplate(w, "layout.html", LoginTemplateData{Email: form.Email, Error: "Invalid email or password"})
			return
		}

		session, err := ctx.SessionStore.New(r, "session")
		if err != nil {
			http.Error(w, "Session error", http.StatusInternalServerError)
			return
		}

		session.Values["user_id"] = user.ID
		err = session.Save(r, w)
		if err != nil {
			println(err.Error())
			http.Error(w, "Session Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Hx-Redirect", "/passages")
		w.WriteHeader(http.StatusNoContent)
	}).Methods("Post")
}
