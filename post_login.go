package main

import (
	"context"
	"errors"
	"main/view"
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
                view.App(view.AppModel{
                    Page: view.LoginPageModel{
                        Email: form.Email,
                        Error: "Invalid email or password",
                    },
                }).Render(r.Context(), w)
			} else {
				println(err.Error())
				http.Error(w, "Database Error", http.StatusInternalServerError)
			}
			return
		}

		if user.Password != form.Password {
            view.App(view.AppModel{
                Page: view.LoginPageModel{
                    Email: form.Email,
                    Error: "Invalid email or password",
                },
            }).Render(r.Context(), w)
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
