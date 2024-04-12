package main

import (
	"main/domain_model"
	"main/view"
	"net/http"
)

func (ctx ServerContext) userRoutes(router *http.ServeMux) {
	router.Handle("GET /register", AuthMiddleware(false, HandleErrors(func(w http.ResponseWriter, r *http.Request) error {
        page := view.RegisterPage(view.RegisterPageModel{})
        return ctx.RenderPage(w, r, page)
	})))

	router.Handle("POST /register", AuthMiddleware(false, HandleErrors(func(w http.ResponseWriter, r *http.Request) error {
        firstName := r.FormValue("first_name")
        lastName := r.FormValue("last_name")
        email := r.FormValue("email_name")
        password := r.FormValue("password")

        name, err := domain_model.NewUserName(firstName, lastName)
        if err != nil {
            return err
        }
        emailAddress, err := domain_model.NewUserEmail(email)
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

    router.Handle("GET /login", AuthMiddleware(false, HandleErrors(func(w http.ResponseWriter, r *http.Request) error {
        page := view.LoginPage(view.LoginPageModel{})
        return ctx.RenderPage(w, r, page)
	})))

	router.Handle("POST /login", AuthMiddleware(false, HandleErrors(func(w http.ResponseWriter, r *http.Request) error {
        email := r.FormValue("email")
        password := r.FormValue("password")

        user, err := ctx.UserRepo.GetByEmail(email)
        if err != nil {
            return err
        }

        if user == nil || !user.ValidatePassword(password) {
            viewModel := view.LoginPageModel{
                Email: email,
                Error: "Invalid email or password.",
            }
            page := view.LoginPage(viewModel)
            return page.Render(r.Context(), w)
        }

		if _, err = ctx.SessionManager.LogIn(w, r, user.Id()); err != nil {
			return err
		}

		w.Header().Set("Hx-Redirect", "/passages")
		w.WriteHeader(http.StatusNoContent)

		return nil
	})))

	router.Handle("POST /logout", AuthMiddleware(true, HandleErrors(func(w http.ResponseWriter, r *http.Request) error {
        if _, err := ctx.SessionManager.LogOut(w, r); err != nil {
            return err
        }

		w.Header().Set("Hx-Redirect", "/")
		w.WriteHeader(http.StatusNoContent)

        return nil
	})))
}
