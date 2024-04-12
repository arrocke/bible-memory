package main

import (
	"main/view"
	"net/http"
)

func indexRoutes(router *http.ServeMux, ctx *ServerContext) {
	router.Handle("GET /", AuthMiddleware(false, HandleErrors(func(w http.ResponseWriter, r *http.Request) error {
        page := view.PublicIndexPage()

        if r.Header.Get("Hx-Request") == "true" {
            return page.Render(r.Context(), w)
        } else {
            app := view.Page(view.AppModel{}, page)
            return app.Render(r.Context(), w)
        }
	})))
}
