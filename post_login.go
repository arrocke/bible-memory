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

	router.Handle("/users/login", AuthMiddleware(false, HandleErrors(func(w http.ResponseWriter, r *http.Request) error {
		form := LoginForm{
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
		}

		query := `SELECT id, password FROM "user" WHERE email = $1`
		rows, _ := ctx.Conn.Query(context.Background(), query, form.Email)
		user, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[DbUser])
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
                return view.CreateViewEngine(ctx.Conn, r.Context(), w).RenderLoginError(
                    "Invalid email or password",
                    form.Email,
                )
			} else {
                return err
			}
		}

		if user.Password != form.Password {
            return view.CreateViewEngine(ctx.Conn, r.Context(), w).RenderLoginError(
                "Invalid email or password",
                form.Email,
            )
		}

        if _, err = ctx.SessionManager.LogIn(w, r, int(user.ID)); err != nil {
            return err
        }

		w.Header().Set("Hx-Redirect", "/passages")
		w.WriteHeader(http.StatusNoContent)

        return nil
	}))).Methods("Post")
}
