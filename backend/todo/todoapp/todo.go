package todoapp

import (
	"context"

	"github.com/mukhtarkv/workspace/kit/errors"
	"github.com/mukhtarkv/workspace/kit/id"
	"go.opentelemetry.io/otel"
)

type ToDoStorage interface {
	Fetch(ctx context.Context, id string) (*ToDoItem, error)
	List(ctx context.Context) ([]ToDoItem, error)
	Create(ctx context.Context, item *ToDoItem) error
	Update(ctx context.Context, item *ToDoItem) error
	Delete(ctx context.Context, id string) error
}

type ToDoItem struct {
	Id string
	Title string
	Details string
}

// ToDoService represent the service that manage todo
// without implementation detail, this aim to provide API that will be exposed to the handler (HTTP, GRPC)
type ToDoService struct {
	// storage is the representation about how we store a todo
	// this does not expose the implementation detail.
	// the backend implementation could be in memory, DB or other.
	storage ToDoStorage
}

// New create a new ToDoItem service.
func New(storage ToDoStorage) *ToDoService {
	return &ToDoService{
		storage: storage,
	}
}

// Create creates a new todo.
// If the todo already exist in the storage, ErrToDoItemAlreadyExist will be returned.
func (u *ToDoService) Create(ctx context.Context, usr *ToDoItem) error {
	const op = "todo.create"
	ctx, span := otel.Tracer("").Start(ctx, op)
	defer span.End()

	// Assign to the todo an ID.
	usr.Id = id.New()

	err := u.storage.Create(ctx, usr)
	if err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}

// List lists todos.
func (u *ToDoService) List(ctx context.Context) ([]ToDoItem, error) {
	const op = "todo.list"
	ctx, span := otel.Tracer("").Start(ctx, op)
	defer span.End()

	items, err := u.storage.List(ctx)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return items, nil
}

// Delete deletes the todo form the todo storage.
// If the todo does not exist, ErrToDoItemNotFound will be returned.
func (u *ToDoService) Update(ctx context.Context, updatedUsr *ToDoItem, mask_paths []string) error {
	const op = "todo.update"
	ctx, span := otel.Tracer("").Start(ctx, op)
	defer span.End()

	usr, err := u.storage.Fetch(ctx, updatedUsr.Id)
	if err != nil {
		return errors.Wrap(err, op)
	}

	for _, mask := range mask_paths {
		switch mask {
		case "title":
			usr.Title = updatedUsr.Title
		case "details":
			usr.Details = updatedUsr.Details
		}
	}

	err = u.storage.Update(ctx, usr)
	if err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}

// Delete deletes the todo form the todo storage.
func (u *ToDoService) Delete(ctx context.Context, id string) error {
	const op = "todo.delete"
	ctx, span := otel.Tracer("").Start(ctx, op)
	defer span.End()

	err := u.storage.Delete(ctx, id)
	if err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}
