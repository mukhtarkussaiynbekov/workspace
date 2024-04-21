package todoapp

import "context"

type ToDoStorage interface {
	Fetch(ctx context.Context, id string) (ToDoItem, error)
	List(ctx context.Context) ([]ToDoItem, error)
	Create(ctx context.Context, item ToDoItem) error
	Delete(ctx context.Context, id string) error
}

type ToDoItem struct {
	Id string
	Title string
	Details string
}