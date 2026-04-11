package driver

import (
	"context"
	"database/sql/driver"

	"github.com/mainlycricket/starks/internal/messages/query"
)

func (c *Conn) QueryContext(ctx context.Context, q string, args []driver.NamedValue) (driver.Rows, error) {
	if err := query.WriteSimpleQuery(ctx, c.connection, q); err != nil {
		return nil, err
	}

	res, err := query.ReadSimpleQueryResp(ctx, c.connection)
	if err != nil {
		return nil, err
	}

	rows := &Rows{
		rowDescriptionBytes: res.RowDescriptionBytes,
		dataRowBytes:        res.DataRowBytes,
	}

	return rows, nil
}
