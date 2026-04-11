package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mainlycricket/starks/internal/driver"
)

func main() {
	db, err := sql.Open("starks", "host=localhost,port=5432,user=postgres,password=1456,database=postgres")
	if err != nil {
		fmt.Println(err)
	}

	if err := db.Ping(); err != nil {
		fmt.Println("err in ping", err)
	}
}
