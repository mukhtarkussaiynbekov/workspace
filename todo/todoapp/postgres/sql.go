package postgres

import (
	"database/sql"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/mukhtarkv/workspace/kit/errors"
	"github.com/mukhtarkv/workspace/todo/todoapp"
)

// Wrap returns an error annotating err with a stack trace
// at the point Wrap is called, and the supplied message.
// If err is nil, Wrap returns nil.
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}

	var modifiedErr todoapp.Error = todoapp.Error(err.Error())
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			modifiedErr = todoapp.ErrToDoItemAlreadyExist
		default:
			modifiedErr = todoapp.Error(err.Error())
		}
	}

	if errors.Is(err, sql.ErrNoRows) {
		modifiedErr = todoapp.ErrToDoItemNotFound
	}

	return fmt.Errorf("%s: %w", message, errors.New(modifiedErr.Error()))
}