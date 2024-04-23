package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	db "github.com/mukhtarkv/workspace/kit/sql"
	"github.com/mukhtarkv/workspace/todo/todoapp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// Verify interface compliance
var _ todoapp.ToDoStorage = (*ToDoStorage)(nil)

type ToDoStorage struct {
	DB *sqlx.DB
}

func New() (*ToDoStorage, error) {
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", 
			os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), 
			os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"), 
			os.Getenv("POSTGRES_DB"))
	rootdb, err := db.Open(dbUrl)
	if err != nil {
		return nil, err
	}
	return &ToDoStorage{
		DB: rootdb,
	}, nil
}

type todo struct {
	Id string `db:"id"`
	Title string `db:"title"`
	Details string `db:"details"`
}

func (s *ToDoStorage) Fetch(ctx context.Context, id string) (*todoapp.ToDoItem, error) {
	_, span := otel.Tracer("").Start(ctx, "postgres.fetch")
	defer span.End()
	
	q := `SELECT id, title, details
		FROM todo
		WHERE id = $1`
	var res todo
	if err := s.DB.QueryRowxContext(ctx, q, id).StructScan(&res); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "fetching")
		span.SetAttributes(attribute.String("id", id))
		return nil, Wrap(err, "fetching")
	}
	return &todoapp.ToDoItem{
		Id: res.Id,
		Title: res.Title,
		Details: res.Details,
	}, nil
}

func (s *ToDoStorage) List(ctx context.Context) ([]todoapp.ToDoItem, error) {
	_, span := otel.Tracer("").Start(ctx, "postgres.list")
	defer span.End()

	q := `SELECT id, title, details
		FROM todo`
	var entities []todo
	if err := s.DB.SelectContext(ctx, &entities, q); err != nil {
		if err == sql.ErrNoRows {
			return []todoapp.ToDoItem{}, nil
		}
		span.RecordError(err)
		span.SetStatus(codes.Error, "listing todo entities")
		return []todoapp.ToDoItem{}, Wrap(err, "listing todo entities")
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

func (s *ToDoStorage) Create(ctx context.Context, item *todoapp.ToDoItem) error {
	_, span := otel.Tracer("").Start(ctx, "postgres.create")
	defer span.End()

	q := `INSERT INTO todo (id, title, details)
		VALUES (:id, :title, :details)`
	entity := todo{
		Id: item.Id,
		Title: item.Title,
		Details: item.Details,
	}
	if _, err := s.DB.NamedExecContext(ctx, q, entity); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "creating entity")
		span.SetAttributes(
			attribute.String("id", entity.Id),
			attribute.String("title", entity.Title),
			attribute.String("details", entity.Details),
		)
		return Wrap(err, "creating entity")
	}
	return nil
}

func (s *ToDoStorage) Update(ctx context.Context, item *todoapp.ToDoItem) error {
	_, span := otel.Tracer("").Start(ctx, "postgres.update")
	defer span.End()

	q := `UPDATE todo
		SET title = :title,
			details = :details
		WHERE id = :id`
	entity := todo{
		Id: item.Id,
		Title: item.Title,
		Details: item.Details,
	}
	if _, err := s.DB.NamedExecContext(ctx, q, entity); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "updating entity")
		span.SetAttributes(
			attribute.String("id", entity.Id),
			attribute.String("title", entity.Title),
			attribute.String("details", entity.Details),
		)
		return Wrap(err, "updating entity")
	}
	return nil
}

func (s *ToDoStorage) Delete(ctx context.Context, id string) error {
	_, span := otel.Tracer("").Start(ctx, "postgres.delete")
	defer span.End()

	q := `DELETE FROM todo WHERE id = $1`
	if _, err := s.DB.ExecContext(ctx, q, id); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "deleting entity")
		span.SetAttributes(attribute.String("id", id))
		return Wrap(err, "deleting entity")
	}
	return nil
}
