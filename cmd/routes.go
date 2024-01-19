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

	mux.Get("/", app.handleGetTasks)
	mux.Post("/tasks", app.handleCreateTask)
	mux.Put("/tasks/{id}/toggle", app.toggleTask)

	return mux
}
