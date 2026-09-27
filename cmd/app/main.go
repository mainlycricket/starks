package main

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/mainlycricket/starks/internal/driver"
)

var DB *sql.DB

func main() {
	db, err := sql.Open("starks", "host=localhost,port=5432,user=postgres,password=1456,database=asoiaf")
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

	doReadTests()

	doWriteTests()
}

func doReadTests() {
	if err := simpleSelectQuery(); err != nil {
		log.Fatalf("error while simple select query: %v", err)
	}

	if err := contextSelectQuery(); err != nil {
		log.Fatalf("error while context select query: %v", err)
	}

	if err := simpleSelectRowQuery(); err != nil {
		log.Fatalf("error while simple select row query: %v", err)
	}

	if err := contextSelectRowQuery(); err != nil {
		log.Fatalf("error while context select row query: %v", err)
	}
}

func doWriteTests() {
	if err := simpleInsertExec(); err != nil {
		log.Fatalf("error while simple insert exec: %v", err)
	}

	if err := contextInsertExec(); err != nil {
		log.Fatalf("error while context insert exec: %v", err)
	}

	if err := simpleUpdateExec(); err != nil {
		log.Fatalf("error while simple update exec: %v", err)
	}

	if err := contextUpdateExec(); err != nil {
		log.Fatalf("error while context update exec: %v", err)
	}

	if err := simpleDeleteExec(); err != nil {
		log.Fatalf("error while simple delete exec: %v", err)
	}

	if err := contextDeleteExec(); err != nil {
		log.Fatalf("error while context delete exec: %v", err)
	}
}
