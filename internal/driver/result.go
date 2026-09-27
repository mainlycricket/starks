package driver

import (
	"bytes"
	"strconv"
)

type Result struct {
	commandTag   string
	rowsAffected int64
}

// PostgreSQL doesn't return the last inserted id
func (r *Result) LastInsertId() (int64, error) {
	return 0, nil
}

func (r *Result) RowsAffected() (int64, error) {
	return r.rowsAffected, nil
}

func newResult(commandCompleteData []byte) *Result {
	commandCompleteData = bytes.TrimSuffix(commandCompleteData, []byte{0})
	parts := bytes.Split(commandCompleteData, []byte{32})
	tag := string(parts[0])

	var valueString string

	switch tag {
	case "INSERT":
		valueString = string(parts[2])
	default:
		valueString = string(parts[1])
	}

	rowsAffected, _ := strconv.ParseInt(valueString, 10, 64)

	return &Result{
		commandTag:   tag,
		rowsAffected: rowsAffected,
	}
}
