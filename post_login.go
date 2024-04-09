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

	router.HandleFunc("/users/login", HandleErrors(func(w http.ResponseWriter, r *http.Request) error {
		form := LoginForm{
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
		}

		query := `SELECT id, password FROM "user" WHERE email = $1`
		rows, _ := ctx.Conn.Query(context.Background(), query, form.Email)
		user, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[DbUser])
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
                if err := view.CreateViewEngine(ctx.Conn, r.Context(), w).RenderLoginError(
                    "Invalid email or password",
                    form.Email,
                ); err != nil {
                    return err
                }
                return nil
			} else {
                return err
			}
		}

		if user.Password != form.Password {
            if err := view.CreateViewEngine(ctx.Conn, r.Context(), w).RenderLoginError(
                "Invalid email or password",
                form.Email,
            ); err != nil {
                return err
            }
			return nil
		}

		session, err := ctx.SessionStore.New(r, "session")
		if err != nil {
			return err
		}

		session.Values["user_id"] = user.ID
        if err := session.Save(r, w); err != nil {
			return err
		}

		w.Header().Set("Hx-Redirect", "/passages")
		w.WriteHeader(http.StatusNoContent)

        return nil
	})).Methods("Post")
}
