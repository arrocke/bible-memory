package main

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

func PostRegister(router *mux.Router, ctx *ServerContext) {
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

		var id int32
		query := `INSERT INTO "user" (email, first_name, last_name, password) VALUES ($1, $2, $3, $4) RETURNING id`
		err := ctx.Conn.QueryRow(context.Background(), query, form.Email, form.FirstName, form.LastName, form.Password).Scan(&id)
		if err != nil {
			println(err.Error())
			http.Error(w, "Database Error", http.StatusInternalServerError)
			return
		}

		session, err := ctx.SessionStore.New(r, "session")
		if err != nil {
			http.Error(w, "Session error", http.StatusInternalServerError)
			return
		}

		session.Values["user_id"] = id
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, "Session Error", http.StatusInternalServerError)
		}

		w.Header().Set("Hx-Redirect", "/passages")
		w.WriteHeader(http.StatusNoContent)
	}).Methods("Post")
}
