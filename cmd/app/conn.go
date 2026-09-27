package main

import (
	"context"
	"fmt"
)

func connTest() error {
	conn, err := DB.Conn(context.Background())
	if err != nil {
		return err
	}

	res, err := conn.QueryContext(context.Background(), "SELECT 1;")
	if err != nil {
		return err
	}

	columnNames, err := res.Columns()
	if err != nil {
		return err
	}

	fmt.Println("column names in conn Test", columnNames)

	return nil
}
