package postgres

import (
	"log"
	"testing"

	db "github.com/anthonycorbacho/workspace/kit/sql"
	"github.com/anthonycorbacho/workspace/todo/todoapp/assets"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type todoTestSuite struct {
	suite.Suite
	st *storage
}

func TestToDoSuite(t *testing.T) {
	suite.Run(t, new(todoTestSuite))
}

func (s *todoTestSuite) SetupSuite() {
	log.Println("setup suite")

	rootdb, err := db.Open("postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		s.T().Fatalf("opening db connection %v", err)
	}
	s.st = &storage{
		db: rootdb,
	}

	// Apply migration scripts
	if err := db.Migrate(rootdb, "test_todo", assets.SF); err != nil {
		s.T().Fatalf("applying schema migration scripts %v", err)
	}

	// Add test data
	if _, err := sqlx.LoadFile(rootdb, "./testdata/todo.sql"); err != nil {
		s.T().Fatalf("loading test data %v", err)
	}
}

func (s *todoTestSuite) TearDownSuite() {
	defer s.st.db.Close()
}

func (s *todoTestSuite) TestFetch() {
	todo, err := s.st.Fetch("xid1")
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), todo)
	assert.Equal(s.T(), "xid1", todo.id)
	assert.Equal(s.T(), "buy grocery", todo.title)
	assert.Equal(s.T(), "buy milk and bread", todo.details)
}
