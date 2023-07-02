package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDb *sql.DB

func TestMain(m *testing.M) {
	conn, err := sql.Open("postgres", "postgres://postgres:123456@localhost/go-bank?sslmode=disable")

	if err != nil {
		log.Fatal("Cant connect to database: ", err)
	}

	testDb = conn
	testQueries = New(testDb)

	os.Exit(m.Run())
}
