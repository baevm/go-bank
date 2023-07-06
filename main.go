package main

import (
	"database/sql"
	"go-bank/api"
	"go-bank/config"
	"log"

	db "go-bank/db/sqlc"

	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Load(".")

	if err != nil {
		log.Fatal("Cant read config file: ", err)
	}

	conn, err := sql.Open("postgres", cfg.DB_DSN)

	if err != nil {
		log.Fatal("Cant connect to database: ", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(cfg, store)

	if err != nil {
		log.Fatal("Cant start server: ", err)
	}

	err = server.Start(cfg.SRV_ADDR)

	if err != nil {
		log.Fatal("Cant start server: ", err)
	}
}
