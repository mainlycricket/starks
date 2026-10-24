package driver

import (
	"database/sql/driver"
	"io"

	byte_utils "github.com/mainlycricket/starks/pkg/byte-utils"
	"github.com/mainlycricket/starks/pkg/operators"
)

type Rows struct {
	resultSetIdx int
	resultSets   []RowsResultSet
	rsFlag       bool
}

type RowsResultSet struct {
	rowDescriptionBytes []byte
	dataRowBytes        [][]byte
	rowIdx              int
}

func (r *Rows) Columns() []string {
	resultSet := &r.resultSets[r.resultSetIdx]

	count, _ := byte_utils.ParseInt16Bytes(resultSet.rowDescriptionBytes[:2])
	columns := make([]string, 0, count)

	var columnName []byte
	for idx := 2; idx < len(resultSet.rowDescriptionBytes); {
		b := resultSet.rowDescriptionBytes[idx]
		if b == 0 {
			columns = append(columns, string(columnName))
			idx += 19
			columnName = nil
		} else {
			columnName = append(columnName, b)
			idx++
		}
	}

	return columns
}

func (r *Rows) Close() error {
	return nil
}

func (r *Rows) HasNextResultSet() bool {
	return len(r.resultSets) > r.resultSetIdx
}

func (r *Rows) NextResultSet() error {
	if r.HasNextResultSet() {
		return nil
	}

	return io.EOF
}

func (r *Rows) Next(dest []driver.Value) error {
	resultSet := &r.resultSets[r.resultSetIdx]

	if resultSet.rowIdx >= len(resultSet.dataRowBytes) {
		r.resultSetIdx++
		return io.EOF
	}

	rowBytes := resultSet.dataRowBytes[resultSet.rowIdx]
	var bytesCounter int

	count, _ := byte_utils.ParseInt16Bytes(rowBytes[:2])
	bytesCounter += 2

	for i := range count {
		dataLen, _ := byte_utils.ParseInt32Bytes(rowBytes[bytesCounter : bytesCounter+4])
		bytesCounter += 4
		if dataLen > 0 {
			dest[i] = rowBytes[bytesCounter : bytesCounter+int(dataLen)]
		} else {
			dest[i] = nil
		}
		bytesCounter += operators.Ternary(dataLen > 0, int(dataLen), 0)
	}

	resultSet.rowIdx++

	return nil
}
