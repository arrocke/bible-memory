package main

import (
	"main/view"
	"net/http"

	"github.com/gorilla/mux"
)

func GetCreatePassage(router *mux.Router, ctx *ServerContext) {
	router.HandleFunc("/passages/new", func(w http.ResponseWriter, r *http.Request) {
		session, err := GetSession(r, ctx)
		if err != nil {
			http.Error(w, "Session Error", http.StatusInternalServerError)
			return
		}
		if session == nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		if r.Header.Get("Hx-Current-Url") == "" {
            model, err := LoadPassagesPageModel(ctx.Conn, *session.user_id, GetClientDate(r), view.AddPassagePageModel{})
            if err != nil {
                http.Error(w, "Database Error", http.StatusInternalServerError)
                return
            }

            view.App(model).Render(r.Context(), w)
		} else {
            view.AddPassagePage(view.AddPassagePageModel{}).Render(r.Context(), w)
		}

	}).Methods("Get")
}
