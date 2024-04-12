package main

import (
	"main/domain_model"
	"main/services"
	"main/view"
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
            if err, ok := err.(domain_model.DomainError); ok {
                message := "Unknown Error"
                switch err.Object {
                case "UserName":
                    switch err.Code {
                    case "FirstNameEmpty":
                        message = "First name is required."
                    case "LastNameEmpty":
                        message = "Last name is required."
                    }
                case "UserEmail":
                    switch err.Code {
                    case "Empty":
                        message = "Email is required."
                    case "Format":
                        message = "Email format is invalid."
                    }
                }
                
                engine := view.CreateViewEngine(ctx.Conn, r.Context(), w)
                return engine.RenderProfileForm(userId, message)
            } else {
			    return err
            }
		}

		w.Header().Set("Hx-Redirect", "/passages")
		w.WriteHeader(http.StatusNoContent)

		return nil
	}))).Methods("Put")
}
