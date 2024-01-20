package repository

import (
	"context"
	"database/sql"

	"github.com/Morpa/htmx-go-todo/internal/models"
)

type DatabaseRepo interface {
	Connection() *sql.DB
	FetchTasks() ([]*models.Item, error)
	FetchCount() (int, error)
	FetchCompletedCount() (int, error)
	ToggleTask(ID int) (*models.Item, error)
	DeleteTask(ctx context.Context, ID int) error
	FetchTask(ID int) (*models.Item, error)
	UpdateTask(ID int, title string) (*models.Item, error)
	OrderTask(ctx context.Context, values []int) error
	InsertTask(title string) (*models.Item, error)
}
