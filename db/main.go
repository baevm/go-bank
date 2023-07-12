package db

import (
	"database/sql"
	"fmt"
	sqlc "go-bank/db/sqlc"
	"log"

	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func Start(dbDsn string) (sqlc.Store, error) {
	conn, err := sql.Open("postgres", dbDsn)

	if err != nil {
		return nil, fmt.Errorf("cant connect to database: %s", err)
	}

	m, err := migrate.New("file://db/migration", dbDsn)

	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %s", err)
	}

	err = m.Up()

	if err != nil && err != migrate.ErrNoChange {
		return nil, fmt.Errorf("failed to run migrate: %s", err)
	}

	log.Println("Database migrated succesfully")
	db := sqlc.NewStore(conn)

	return db, nil
}
