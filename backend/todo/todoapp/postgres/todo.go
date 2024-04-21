package postgres

import "github.com/jmoiron/sqlx"

type storage struct {
	db *sqlx.DB
}

type todo struct {
	id string `db:"id"`
	title string `db:"title"`
	details string `db:"details"`
}

func (s *storage) Fetch(id string) (*todo, error) {
	return nil, nil
}