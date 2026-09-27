package main

import (
	"context"
	"fmt"
)

func simpleQuery() error {
	res, err := DB.Query("SELECT 1 as data, 2 as data2;")
	if err != nil {
		return err
	}

	if res.Err() != nil {
		return err
	}

	cols, err := res.Columns()
	if err != nil {
		return err
	}

	fmt.Println("column names in simple query: ", cols)

	for res.Next() {
		var (
			d  int64
			d2 *int64
		)
		if err := res.Scan(&d, &d2); err != nil {
			return err
		}
		fmt.Println("data", d, *d2)
	}

	return nil
}

func contextQuery() error {
	res, err := DB.QueryContext(context.Background(), "SELECT 1 as data, 2 as data2;")
	if err != nil {
		return err
	}

	if res.Err() != nil {
		return err
	}

	cols, err := res.Columns()
	if err != nil {
		return err
	}

	fmt.Println("column names in query context: ", cols)

	for res.Next() {
		var (
			d  int64
			d2 *int64
		)
		if err := res.Scan(&d, &d2); err != nil {
			return err
		}
		fmt.Println("data", d, *d2)
	}

	return nil
}
