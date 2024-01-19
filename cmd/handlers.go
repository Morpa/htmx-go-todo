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

	completedCount, err := app.DB.FetchCompletedCount()
	if err != nil {
		log.Printf("error fetching completed count: %v", err)
		return
	}

	data := models.Tasks{
		Items:          tasks,
		Count:          count,
		CompletedCount: completedCount,
	}

	tmpl.ExecuteTemplate(w, "Base", data)
}

func (app *application) HandleCreateTask(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	if title == "" {
		return
	}

	_, err := app.DB.InsertTask(title)
	if err != nil {
		log.Printf("error insert task: %v", err)
		return
	}

	_, err = app.DB.FetchCount()
	if err != nil {
		log.Printf("error fetching count: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	tmpl.ExecuteTemplate(w, "Form", nil)
}
