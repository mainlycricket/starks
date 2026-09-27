package main

import (
	"context"
	"fmt"
)

func simpleDeleteExec() error {
	query := "DELETE FROM houses WHERE id = $1;"

	res, err := DB.Exec(query, 1)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	fmt.Println("rows affected in simple delete exec", rowsAffected)

	return nil
}

func contextDeleteExec() error {
	query := "DELETE FROM houses WHERE id = $1;"

	res, err := DB.ExecContext(context.Background(), query, 2)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	fmt.Println("rows affected in context delete exec", rowsAffected)

	return nil
}
