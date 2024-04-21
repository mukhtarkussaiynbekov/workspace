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

func (s *storage) Fetch(ctx context.Context, id string) (todoapp.ToDoItem, error) {
	q := `SELECT id, title, details
		FROM todo
		WHERE id = $1`
	var res todo
	if err := s.db.QueryRowxContext(ctx, q, id).StructScan(&res); err != nil {
		return todoapp.ToDoItem{}, err
	}
	return todoapp.ToDoItem{
		Id: res.Id,
		Title: res.Title,
		Details: res.Details,
	}, nil
}

func (s *storage) List(ctx context.Context) ([]todoapp.ToDoItem, error) {
	q := `SELECT id, title, details
		FROM todo`
	var entities []todo
	if err := s.db.SelectContext(ctx, &entities, q); err != nil {
		return []todoapp.ToDoItem{}, err
	}
	var res []todoapp.ToDoItem
	for _, entity := range entities {
		res = append(res, todoapp.ToDoItem{
			Id: entity.Id,
			Title: entity.Title,
			Details: entity.Details,
		})
	}
	return res, nil
}
