package main

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/mainlycricket/starks/internal/driver"
)

var DB *sql.DB

func main() {
	db, err := sql.Open("starks", "host=localhost,port=5432,user=postgres,password=1456,database=postgres")
	if err != nil {
		log.Fatalf("error while opening connection: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("err in ping: %v", err)
	}

	if err := db.PingContext(context.Background()); err != nil {
		log.Fatalf("err in ping: %v", err)
	}

	DB = db

	if err := connTest(); err != nil {
		log.Fatalf("error in Conn Test: %v", err)
	}

	if err := simpleQuery(); err != nil {
		log.Fatalf("error while simple query: %v", err)
	}

	if err := contextQuery(); err != nil {
		log.Fatalf("error while context query: %v", err)
	}
}
