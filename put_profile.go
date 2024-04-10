package main

import (
	"main/services"
	"net/http"

	"github.com/gorilla/mux"
)

func PutProfile(router *mux.Router, ctx *ServerContext) {
	router.Handle("/users/profile", AuthMiddleware(true, HandleErrors(func(w http.ResponseWriter, r *http.Request) error {
		userId := GetUserId(r)

		if err := ctx.UserService.UpdateProfile(services.UpdateProfileRequest{
			Id:           userId,
			EmailAddress: r.FormValue("email"),
			FirstName:    r.FormValue("first_name"),
			LastName:     r.FormValue("last_name"),
			Password:     r.FormValue("password"),
		}); err != nil {
			return err
		}

		w.Header().Set("Hx-Redirect", "/passages")
		w.WriteHeader(http.StatusNoContent)

		return nil
	}))).Methods("Put")
}
