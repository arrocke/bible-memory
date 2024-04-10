package main

import (
	"main/services"
	"net/http"

	"github.com/gorilla/mux"
)

func PostRegister(router *mux.Router, ctx *ServerContext) {
	router.Handle("/users/register", AuthMiddleware(false, HandleErrors(func(w http.ResponseWriter, r *http.Request) error {
		id, err := ctx.UserService.Create(services.CreateUserRequest{
			EmailAddress: r.FormValue("email"),
			FirstName:    r.FormValue("first_name"),
			LastName:     r.FormValue("last_name"),
			Password:     r.FormValue("password"),
		})
		if err != nil {
			return err
		}

		if _, err := ctx.SessionManager.LogIn(w, r, id); err != nil {
			return err
		}

		w.Header().Set("Hx-Redirect", "/passages")
		w.WriteHeader(http.StatusNoContent)

		return nil
	}))).Methods("Post")
}
