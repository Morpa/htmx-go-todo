package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/Morpa/htmx-go-todo/internal/db"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	err := db.OpenDB()
	if err != nil {
		log.Panic(err)
	}
	defer db.CloseDB()

	err = db.SetupDB()
	if err != nil {
		log.Panic(err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	r.Get("/", func(w http.ResponseWriter, _ *http.Request) {
		tmpl, _ := template.New("").ParseFiles("templates/index.html")
		tmpl.ExecuteTemplate(w, "Base", nil)
	})
	http.ListenAndServe(":3000", r)
}
