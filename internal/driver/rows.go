package driver

import (
	"database/sql/driver"
	"errors"

	byte_utils "github.com/mainlycricket/starks/pkg/byte-utils"
	"github.com/mainlycricket/starks/pkg/operators"
)

type Rows struct {
	rowDescriptionBytes []byte
	dataRowBytes        [][]byte
	rowIdx              int
}

func (r *Rows) Columns() []string {
	count, _ := byte_utils.ParseInt16Bytes(r.rowDescriptionBytes[:2])
	columns := make([]string, 0, count)

	var columnName []byte
	for idx := 2; idx < len(r.rowDescriptionBytes); {
		b := r.rowDescriptionBytes[idx]
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

func (r *Rows) Next(dest []driver.Value) error {
	if len(r.dataRowBytes) == r.rowIdx {
		return errors.New("end of rows")
	}

	rowBytes := r.dataRowBytes[r.rowIdx]
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

	r.rowIdx++

	return nil
}
