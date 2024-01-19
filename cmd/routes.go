package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Logger)

	mux.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// mux.Get("/", func(w http.ResponseWriter, _ *http.Request) {
	// 	tmpl.ExecuteTemplate(w, "Base", nil)
	// })
	mux.Get("/", app.HandleGetTasks)

	return mux
}
