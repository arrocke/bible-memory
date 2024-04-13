package main

import (
	"main/domain_model"
	"main/view"
	"net/http"

	"github.com/jackc/pgx/v5"
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

    
	router.Handle("GET /profile", AuthMiddleware(true, HandleErrors(func(w http.ResponseWriter, r *http.Request) error {
        userId := GetUserId(r)

        query := `SELECT email, first_name, last_name FROM "user" WHERE id = $1`
        rows, _ := ctx.Conn.Query(r.Context(), query, userId)
        model, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[view.ProfilePageModel])
        if err != nil {
            return err
        }

        page := view.ProfilePage(model)
        return ctx.RenderPage(w, r, page)
	})))

    router.Handle("PUT /profile", AuthMiddleware(true, HandleErrors(func(w http.ResponseWriter, r *http.Request) error {
		userId := GetUserId(r)

        user, err := ctx.UserRepo.Get(userId)
        if err != nil {
            return err
        }

        if user == nil {
            // TODO: handle not found
        }

        name, err := domain_model.NewUserName(r.FormValue("first_name"), r.FormValue("last_name"))
        if err != nil {
            return err
        }
        user.ChangeName(name)

        email, err := domain_model.NewUserEmail(r.FormValue("email"))
        if err != nil {
            return err
        }
        user.ChangeEmail(email)

        
        password := r.FormValue("password")
        if (password != "") {
            user.ChangePassword(password)
        }

        if err := ctx.UserRepo.Commit(*user); err != nil {
            return err
        }

		w.Header().Set("Hx-Redirect", "/passages")
		w.WriteHeader(http.StatusNoContent)

		return nil
	})))
}
