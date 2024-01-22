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

	router.HandleFunc("/users/login", func(w http.ResponseWriter, r *http.Request) {
		form := LoginForm{
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
		}

		var expectedPassword string
		query := `SELECT password FROM "user" WHERE email = $1`
		err := ctx.Conn.QueryRow(context.Background(), query, form.Email).Scan(&expectedPassword)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				http.Error(w, "Wrong email or password", http.StatusUnauthorized)
			} else {
				http.Error(w, "Database Error", http.StatusInternalServerError)
			}
			return
		}

		if expectedPassword != form.Password {
			http.Error(w, "Wrong email or password", http.StatusUnauthorized)
			return
		}

		w.Header().Set("Hx-Redirect", "/passages")
		w.WriteHeader(http.StatusNoContent)
	}).Methods("Post")
}
