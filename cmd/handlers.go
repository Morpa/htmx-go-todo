package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Morpa/htmx-go-todo/internal/models"
	"github.com/go-chi/chi/v5"
)

func (app *application) handleGetTasks(w http.ResponseWriter, r *http.Request) {
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

func (app *application) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	if title == "" {
		tmpl.ExecuteTemplate(w, "Form", nil)
		return
	}

	item, err := app.DB.InsertTask(title)
	if err != nil {
		log.Printf("error insert task: %v", err)
		return
	}

	count, err := app.DB.FetchCount()
	if err != nil {
		log.Printf("error fetching count: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	tmpl.ExecuteTemplate(w, "Form", nil)
	tmpl.ExecuteTemplate(w, "Item", map[string]any{"Item": item, "SwapOOB": true})
	tmpl.ExecuteTemplate(w, "TotalCount", map[string]any{"Count": count, "SwapOOB": true})
}

func (app *application) toggleTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Printf("error parsing id into int: %v", err)
		return
	}

	_, err = app.DB.ToggleTask(id)
	if err != nil {
		log.Printf("error toggling task: %v", err)
		return
	}
	completedCount, err := app.DB.FetchCompletedCount()
	if err != nil {
		log.Printf("error fetching completed count: %v", err)
		return
	}
	tmpl.ExecuteTemplate(w, "CompletedCount", map[string]any{"Count": completedCount, "SwapOOB": true})
}

func (app *application) handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Printf("error parsing id into int: %v", err)
		return
	}

	err = app.DB.DeleteTask(r.Context(), id)
	if err != nil {
		log.Printf("error deleting the task: %v", err)
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
	tmpl.ExecuteTemplate(w, "TotalCount", map[string]any{"Count": count, "SwapOOB": true})
	tmpl.ExecuteTemplate(w, "CompletedCount", map[string]any{"Count": completedCount, "SwapOOB": true})
}

func (app *application) handleEditTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Printf("error parsing id into int: %v", err)
		return
	}

	task, err := app.DB.FetchTask(id)
	if err != nil {
		log.Printf("error fetching task with id: %d %v", id, err)
		return
	}
	tmpl.ExecuteTemplate(w, "Item", map[string]any{"Item": task, "Editing": true})
}

func (app *application) handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Printf("error parsing id into int: %v", err)
		return
	}
	title := r.FormValue("title")
	if title == "" {
		return
	}

	task, err := app.DB.UpdateTask(id, title)
	if err != nil {
		log.Printf("error fetching task with id: %d %v", id, err)
		return
	}
	tmpl.ExecuteTemplate(w, "Item", map[string]any{"Item": task})
}

func (app *application) handleOrderTasks(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Printf("error parsing form: %v", err)
		return
	}

	var values []int
	for k, v := range r.Form {
		if k == "item" {
			for _, v := range v {
				value, err := strconv.Atoi(v)
				if err != nil {
					log.Printf("error parsing id into int: %v", err)
					return
				}
				values = append(values, value)
			}
		}
	}

	err = app.DB.OrderTask(r.Context(), values)
	if err != nil {
		log.Printf("error ordering tasks: %v", err)
		return
	}
}
