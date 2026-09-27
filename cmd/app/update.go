package main

import (
	"context"
	"crypto/rand"
	"fmt"
)

func simpleUpdateExec() error {
	query := "UPDATE houses SET name = $1 WHERE id = $2;"

	res, err := DB.Exec(query, rand.Text(), 1)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	fmt.Println("rows affected in simple update exec", rowsAffected)

	return nil
}

func contextUpdateExec() error {
	query := "UPDATE houses SET name = $1 WHERE id = $2;"

	res, err := DB.ExecContext(context.Background(), query, rand.Text(), 2)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	fmt.Println("rows affected in context update exec", rowsAffected)

	return nil
}
