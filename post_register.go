package main

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

func PostRegister(router *mux.Router, ctx *ServerContext) {
	type registerForm struct {
		Email           string
		FirstName       string
		LastName        string
		Password        string
		ConfirmPassword string
	}

	router.HandleFunc("/users/register", HandleErrors(func(w http.ResponseWriter, r *http.Request) error {
		form := registerForm{
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
			return err
		}

		session, err := ctx.SessionStore.New(r, "session")
		if err != nil {
			return err
		}

		session.Values["user_id"] = id
        if err := session.Save(r, w); err != nil {
			return err
		}

		w.Header().Set("Hx-Redirect", "/passages")
		w.WriteHeader(http.StatusNoContent)

        return nil
	})).Methods("Post")
}
