package dbrepo

import (
	"context"
	"database/sql"

	"github.com/Morpa/htmx-go-todo/internal/models"
)

type SqliteDBRepo struct {
	DB *sql.DB
}

func (m *SqliteDBRepo) Connection() *sql.DB {
	return m.DB
}

func (m *SqliteDBRepo) FetchTasks() ([]*models.Item, error) {
	query := "select id, title, completed from tasks order by position"

	rows, err := m.DB.Query(query)
	if err != nil {
		return []*models.Item{}, err
	}
	defer rows.Close()

	var items []*models.Item

	for rows.Next() {
		var item models.Item
		err := rows.Scan(
			&item.ID,
			&item.Title,
			&item.Completed,
		)
		if err != nil {
			return []*models.Item{}, err
		}

		items = append(items, &item)
	}
	return items, nil
}

func (m *SqliteDBRepo) FetchTask(ID int) (models.Item, error) {
	var item models.Item
	query := "select id, title, completed from tasks where id = (?)"
	row := m.DB.QueryRow(query, ID)
	err := row.Scan(
		&item.ID,
		&item.Title,
		&item.Completed,
	)
	if err != nil {
		return models.Item{}, err
	}

	return item, nil
}

func (m *SqliteDBRepo) UpdateTask(ID int, title string) (models.Item, error) {
	var item models.Item
	query := `update tasks set title = (?) where id = (?) returning id, title, completed`
	row := m.DB.QueryRow(query, title, ID)
	err := row.Scan(
		&item.ID,
		&item.Title,
		&item.Completed,
	)
	if err != nil {
		return models.Item{}, err
	}

	return item, nil
}

func (m *SqliteDBRepo) FetchCount() (int, error) {
	query := "select count(*) from tasks"

	var count int
	row := m.DB.QueryRow(query)
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (m *SqliteDBRepo) FetchCompletedCount() (int, error) {
	query := "select count(*) from tasks where completed = 1"

	var count int
	row := m.DB.QueryRow(query)
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (m *SqliteDBRepo) InsertTask(title string) (models.Item, error) {
	count, err := m.FetchCount()
	if err != nil {
		return models.Item{}, err
	}

	var id int
	query := `insert into tasks (title, position) values (?, ?) returning id`
	row := m.DB.QueryRow(query, title, count)
	err = row.Scan(&id)
	if err != nil {
		return models.Item{}, err
	}

	item := models.Item{
		ID:        id,
		Title:     title,
		Completed: false,
	}

	return item, nil
}

func (m *SqliteDBRepo) DeleteTask(ctx context.Context, ID int) error {
	query := `delete from tasks where id = (?)`
	_, err := m.DB.Exec(query, ID)
	if err != nil {
		return err
	}

	query = `select id from tasks order by position`
	rows, err := m.DB.Query(query)
	if err != nil {
		return err
	}

	var ids []int
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			return err
		}

		ids = append(ids, id)
	}
	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	for idx, id := range ids {
		query := `update tasks set position = (?) where id = (?)`
		_, err := m.DB.Exec(query, idx, id)
		if err != nil {
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (m *SqliteDBRepo) OrderTask(ctx context.Context, values []int) error {
	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	for i, v := range values {
		query := `update tasks set position = (?) where id = (?)`
		_, err := m.DB.Exec(query, i, v)
		if err != nil {
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (m *SqliteDBRepo) ToggleTask(ID int) (*models.Item, error) {
	var item models.Item

	query := "update tasks set completed = case when completed = 1 then 0 else 1 end where id = (?) returning id, title, completed"
	row := m.DB.QueryRow(query, ID)
	err := row.Scan(&item.ID, &item.Title, &item.Completed)
	if err != nil {
		return nil, err
	}
	return &item, nil
}
