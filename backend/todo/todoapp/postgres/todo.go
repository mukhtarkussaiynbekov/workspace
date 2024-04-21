package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/mukhtarkv/workspace/todo/todoapp"
)

type storage struct {
	db *sqlx.DB
}

type todo struct {
	Id string `db:"id"`
	Title string `db:"title"`
	Details string `db:"details"`
}

func (s *storage) Fetch(id string) (todoapp.ToDoItem, error) {
	q := `SELECT id, title, details
		FROM todo
		WHERE id = $1`
	var res todo
	if err := s.db.QueryRowxContext(context.Background(), q, id).StructScan(&res); err != nil {
		return todoapp.ToDoItem{}, err
	}
	return todoapp.ToDoItem{
		Id: res.Id,
		Title: res.Title,
		Details: res.Details,
	}, nil
}