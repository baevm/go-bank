package db

import (
	"database/sql"
	"go-bank/config"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDb *sql.DB

func TestMain(m *testing.M) {
	cfg, err := config.Load("../..")
	if err != nil {
		log.Fatal("Cant read cfg file: ", err)
	}

	conn, err := sql.Open("postgres", cfg.DB_DSN)

	if err != nil {
		log.Fatal("Cant connect to database: ", err)
	}

	testDb = conn
	testQueries = New(testDb)

	os.Exit(m.Run())
}
