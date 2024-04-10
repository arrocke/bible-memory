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
		email := r.FormValue("email")
		password := r.FormValue("password")

		id, err := ctx.UserService.ValidatePassword(email, password)
		if err != nil {
			return err
		}

		if id == nil {
			return view.CreateViewEngine(ctx.Conn, r.Context(), w).RenderLoginError(
				"Invalid email or password",
				email,
			)
		}

		if _, err = ctx.SessionManager.LogIn(w, r, *id); err != nil {
			return err
		}

		w.Header().Set("Hx-Redirect", "/passages")
		w.WriteHeader(http.StatusNoContent)

		return nil
	}))).Methods("Post")
}
