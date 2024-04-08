package main

import (
	"main/view"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetPassageEdit(router *mux.Router, ctx *ServerContext) {
	router.HandleFunc("/passages/{Id}", func(w http.ResponseWriter, r *http.Request) {
		session, err := GetSession(r, ctx)
		if err != nil {
			http.Error(w, "Session Error", http.StatusInternalServerError)
			return
		}
		if session == nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars["Id"], 10, 32)
		if err != nil {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		if r.Header.Get("Hx-Current-Url") == "" {
            view.CreateViewEngine(ctx.Conn, r.Context(), w).RenderPassageEdit(int(*session.user_id), int(id), GetClientDate(r))
		} else {
            view.CreateViewEngine(ctx.Conn, r.Context(), w).RenderPassageEditPartial(int(*session.user_id), int(id))
		}
	}).Methods("Get")
}
