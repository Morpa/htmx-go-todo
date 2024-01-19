package main

import (
	"log"
	"net/http"

	"github.com/Morpa/htmx-go-todo/internal/models"
)

func (app *application) HandleGetTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := app.DB.FetchTasks()
	if err != nil {
		log.Printf("error fetching tasks: %v", err)
		return
	}

	count, err := app.DB.FetchCount()
	if err != nil {
		log.Printf("error fetching count: %v", err)
		return
	}

	data := models.Tasks{
		Items:          tasks,
		Count:          count,
		CompletedCount: 0,
	}

	tmpl.ExecuteTemplate(w, "Base", data)
}
