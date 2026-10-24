package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
)

func doStatementTest() {
	if err := statementExecTest(); err != nil {
		log.Fatalf("error in statementExecTest: %v", err)
	}

	if err := statmentQueryTest(); err != nil {
		log.Fatalf("error in statmentQueryTest: %v", err)
	}
}

func statmentQueryTest() error {
	q := "SELECT id, name FROM houses WHERE id > $1;"
	stmt, err := DB.PrepareContext(context.Background(), q)

	if err != nil {
		return err
	}

	res, err := stmt.QueryContext(context.Background(), 1)
	if err != nil {
		return err
	}
	if err := res.Err(); err != nil {
		return err
	}

	cols, err := res.Columns()
	if err != nil {
		return err
	}

	fmt.Println("column names in statement query:", cols)

	var houses []House

	for res.Next() {
		var house House
		if err := res.Scan(&house.Id, &house.Name); err != nil {
			return err
		}
		houses = append(houses, house)
	}

	fmt.Println("houses list in statement query:", houses)

	return nil
}

func statementExecTest() error {
	q := "INSERT INTO houses (name) VALUES ($1);"
	stmt, err := DB.PrepareContext(context.Background(), q)
	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(context.Background(), rand.Text())
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	fmt.Println("rows affected in statement exec test", rowsAffected)

	return nil
}
