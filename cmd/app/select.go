package main

import (
	"context"
	"fmt"
)

type House struct {
	Id   int64
	Name string
}

func simpleSelectQuery() error {
	query := "SELECT id, name FROM houses WHERE id = $1 OR id = $2;"
	res, err := DB.Query(query, 1, 2)
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

	fmt.Println("column names in simple query:", cols)

	var houses []House

	for res.Next() {
		var house House
		if err := res.Scan(&house.Id, &house.Name); err != nil {
			return err
		}
		houses = append(houses, house)
	}

	fmt.Println("houses list in simple query:", houses)

	return nil
}

func contextSelectQuery() error {
	query := "SELECT id, name FROM houses WHERE id = $1 OR id = $2;"

	res, err := DB.QueryContext(context.Background(), query, 1, 2)
	if err != nil {
		return err
	}

	if err := res.Err(); err != nil {
		return err
	}

	defer res.Close()

	cols, err := res.Columns()
	if err != nil {
		return err
	}

	fmt.Println("column names in query context:", cols)

	var houses []House

	for res.Next() {
		var house House

		if err := res.Scan(&house.Id, &house.Name); err != nil {
			return err
		}

		houses = append(houses, house)
	}

	fmt.Println("houses list in context query:", houses)

	return nil
}

func simpleSelectRowQuery() error {
	query := "SELECT id, name FROM houses WHERE id = $1;"
	row := DB.QueryRow(query, 1)

	if err := row.Err(); err != nil {
		return err
	}

	var house House
	if err := row.Scan(&house.Id, &house.Name); err != nil {
		return err
	}

	fmt.Println("house in simple row query", house)

	return nil
}

func contextSelectRowQuery() error {
	query := "SELECT id, name FROM houses WHERE id = $1;"

	row := DB.QueryRowContext(context.Background(), query, 1)

	var house House
	if err := row.Scan(&house.Id, &house.Name); err != nil {
		return err
	}

	fmt.Println("houses in contextRowQuery:", house)

	return nil
}
