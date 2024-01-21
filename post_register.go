package main

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func PostRegister(router *mux.Router, conn *pgxpool.Pool) {
	type RegisterForm struct {
		Email           string
		FirstName       string
		LastName        string
		Password        string
		ConfirmPassword string
	}

	router.HandleFunc("/users/register", func(w http.ResponseWriter, r *http.Request) {
		form := RegisterForm{
			Email:           r.FormValue("email"),
			FirstName:       r.FormValue("first_name"),
			LastName:        r.FormValue("last_name"),
			Password:        r.FormValue("password"),
			ConfirmPassword: r.FormValue("confirm_password"),
		}

		query := `INSERT INTO "user" (email, first_name, last_name, password) VALUES ($1, $2, $3, $4)`
		_, err := conn.Exec(context.Background(), query, form.Email, form.FirstName, form.LastName, form.Password)
		if err != nil {
			println(err.Error())
			http.Error(w, "Database Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Hx-Redirect", "/passages")
		w.WriteHeader(http.StatusNoContent)
	}).Methods("Post")
}
