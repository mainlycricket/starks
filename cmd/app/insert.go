package main

import (
	"context"
	"crypto/rand"
	"fmt"
)

func simpleInsertExec() error {
	query := "INSERT INTO houses (name) VALUES ($1);"

	res, err := DB.Exec(query, rand.Text())
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	fmt.Println("rows affected in simple insert exec", rowsAffected)

	return nil
}

func contextInsertExec() error {
	query := "INSERT INTO houses (name) VALUES ($1);"

	res, err := DB.ExecContext(context.Background(), query, rand.Text())
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	fmt.Println("rows affected in context insert exec", rowsAffected)

	return nil
}
