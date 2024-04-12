package main

import (
	"main/domain_model"
	"main/view"
	"net/http"
	"strings"
)

func registerRoutes(router *http.ServeMux, ctx *ServerContext) {
	router.Handle("GET /users/register", AuthMiddleware(false, HandleErrors(func(w http.ResponseWriter, r *http.Request) error {
        page := view.RegisterPage(view.RegisterPageModel{})

        if r.Header.Get("Hx-Request") == "true" {
            return page.Render(r.Context(), w)
        } else {
            app := view.Page(view.AppModel{}, page)
            return app.Render(r.Context(), w)
        }
	})))

	router.Handle("POST /users/register", AuthMiddleware(false, HandleErrors(func(w http.ResponseWriter, r *http.Request) error {
        hasError := false
        viewModel := view.RegisterPageModel{
            FirstName: r.FormValue("first_name"),
            LastName: r.FormValue("last_name"),
            Email: r.FormValue("email"),
        }
        password := r.FormValue("password")
        confirmPassword := r.FormValue("confirm_password")

        if viewModel.FirstName == "" {
            viewModel.FirstNameError = "Please enter your first name."
            hasError = true
        }
        if viewModel.LastName == "" {
            viewModel.LastNameError = "Please enter your last name."
            hasError = true
        }
        if viewModel.Email == "" {
            viewModel.EmailError = "Please enter your email."
            hasError = true
        }
        if !strings.Contains(viewModel.Email, "@") {
            viewModel.EmailError = "Please enter a valid email."
            hasError = true
        }
        if password == "" {
            viewModel.PasswordError = "Please enter a password."
            hasError = true
        }
        if confirmPassword != password {
            viewModel.ConfirmPasswordError = "Your passwords do not match."
            hasError = true
        }

        if hasError {
            page := view.RegisterPage(viewModel)
            return page.Render(r.Context(), w)
        }


        name, err := domain_model.NewUserName(viewModel.FirstName, viewModel.LastName)
        if err != nil {
            return err
        }

        emailAddress, err := domain_model.NewUserEmail(viewModel.Email)
        if err != nil {
            return err
        }

        user := domain_model.NewUser(domain_model.NewUserProps{
            Name: name,
            EmailAddress: emailAddress,
            Password: password,
        })

        if err := ctx.UserRepo.Commit(user); err != nil {
            return err
        }

		if _, err := ctx.SessionManager.LogIn(w, r, user.Id()); err != nil {
			return err
		}

		w.Header().Set("Hx-Redirect", "/passages")
		w.WriteHeader(http.StatusNoContent)

		return nil
	})))
}
