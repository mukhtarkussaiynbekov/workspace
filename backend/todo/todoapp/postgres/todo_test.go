package postgres

import (
	"log"
	"testing"

	db "github.com/anthonycorbacho/workspace/kit/sql"
	assets "github.com/anthonycorbacho/workspace/todo/todoapp/assets/db"
	"github.com/jmoiron/sqlx"
)

func setupSuite(t testing.T) func(t testing.T) {
	log.Println("setup suite")

	rootdb, err := db.Open("postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		return func(t testing.T) {
			log.Println("teardown suite")
		}
	}
	defer rootdb.Close()

	// Apply migration scripts
	if err := db.Migrate(rootdb, "todo", assets.SF); err != nil {
		return func(t testing.T) {
			log.Println("teardown suite")
		}
	}
	
	// Add test data
	if _, err := sqlx.LoadFile(rootdb, "./testdata/todo.sql"); err != nil {
		return func(t testing.T) {
			log.Println("teardown suite")
		}
	}

	// Return a function to teardown the test
	return func(t testing.T) {
		log.Println("teardown suite")
	}
}

func TestTodo(t *testing.T) {

}
