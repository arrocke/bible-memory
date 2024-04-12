package main

import (
	"main/view"
	"net/http"
)

func (ctx ServerContext) indexRoutes(router *http.ServeMux) {
	router.Handle("GET /", AuthMiddleware(false, HandleErrors(func(w http.ResponseWriter, r *http.Request) error {
        page := view.PublicIndexPage()
        return ctx.RenderPage(w, r, page)
	})))
}
