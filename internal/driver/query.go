package driver

import (
	"context"
	"database/sql/driver"

	"github.com/mainlycricket/starks/internal/messages/query"
	"github.com/mainlycricket/starks/pkg/logger"
)

func (c *Conn) QueryContext(ctx context.Context, q string, args []driver.NamedValue) (driver.Rows, error) {
	logger.Logger.DebugInfo("Context Query invoked")

	rows, err := c.doExtendedQuery(ctx, q, args)
	if err != nil {
		logger.Logger.DebugError("error in Context Query", err)
	}

	logger.Logger.DebugInfo("Context Query finished")

	return rows, nil
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
