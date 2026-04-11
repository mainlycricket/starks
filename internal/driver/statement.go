package driver

import "database/sql/driver"

type Stmt struct {
}

func (s *Stmt) Close() error {
	return nil
}

func (s *Stmt) NumInput() int {
	return 0
}

func (s *Stmt) Exec(args []driver.Value) (driver.Result, error) {
	return &Result{}, nil
}

func (s *Stmt) Query(args []driver.Value) (driver.Rows, error) {
	return &Rows{}, nil
}
