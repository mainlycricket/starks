package driver

import (
	"context"
	"database/sql/driver"

	"github.com/mainlycricket/starks/internal/messages/query"
)

func (c *Conn) Query(q string, args []driver.NamedValue) (driver.Rows, error) {
	return c.doSimpleQuery(context.Background(), q)
}

func (c *Conn) QueryContext(ctx context.Context, q string, args []driver.NamedValue) (driver.Rows, error) {
	return c.doExtendedQuery(ctx, q, args)
}

func (c *Conn) doExtendedQuery(ctx context.Context, q string, args []driver.NamedValue) (driver.Rows, error) {
	if err := query.WriteExtendedQuery(ctx, c.connection, q, args); err != nil {
		return nil, err
	}

	res, err := query.ReadExtendedQueryResp(ctx, c.connection)
	if err != nil {
		return nil, err
	}

	rows := &Rows{
		rowDescriptionBytes: res.RowDescriptionBytes,
		dataRowBytes:        res.DataRowBytes,
	}

	return rows, nil
}

func (c *Conn) doSimpleQuery(ctx context.Context, q string) (driver.Rows, error) {
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
