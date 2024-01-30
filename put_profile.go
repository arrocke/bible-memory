package main

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

func PutProfile(router *mux.Router, ctx *ServerContext) {
	type ProfileForm struct {
		Email           string
		FirstName       string
		LastName        string
		Password        string
		ConfirmPassword string
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

		form := ProfileForm{
			Email:           r.FormValue("email"),
			FirstName:       r.FormValue("first_name"),
			LastName:        r.FormValue("last_name"),
			Password:        r.FormValue("password"),
			ConfirmPassword: r.FormValue("confirm_password"),
		}

		if form.Password != "" {
			query := `UPDATE "user" SET email = $2, first_name = $3, last_name = $4, password = $5 WHERE id = $1`
			_, err = ctx.Conn.Exec(context.Background(), query, *session.user_id, form.Email, form.FirstName, form.LastName, form.Password)
		} else {
			query := `UPDATE "user" SET email = $2, first_name = $3, last_name = $4 WHERE id = $1`
			_, err = ctx.Conn.Exec(context.Background(), query, *session.user_id, form.Email, form.FirstName, form.LastName)
		}
		if err != nil {
			println(err.Error())
			http.Error(w, "Database Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Hx-Redirect", "/passages")
		w.WriteHeader(http.StatusNoContent)
	}).Methods("Put")
}
