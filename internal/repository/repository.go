package repository

import (
	"context"
	"database/sql"

	"github.com/Morpa/htmx-go-todo/internal/models"
)

type DatabaseRepo interface {
	Connection() *sql.DB
	FetchCount() (int, error)
	FetchTasks() ([]*models.Item, error)
	FetchTask(ID int) (models.Item, error)
	UpdateTask(ID int, title string) (models.Item, error)
	InsertTask(title string) (models.Item, error)
	DeleteTask(ctx context.Context, ID int) error
	OrderTask(ctx context.Context, values []int) error
}
