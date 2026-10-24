package driver

import (
	"context"
	"database/sql/driver"

	"github.com/mainlycricket/starks/internal/messages/query"
	"github.com/mainlycricket/starks/pkg/logger"
)

func (c *Conn) QueryContext(ctx context.Context, q string, args []driver.NamedValue) (driver.Rows, error) {
	logger.Logger.DebugInfo("Context Query invoked")

	var (
		rows driver.Rows
		err  error
	)

	if len(args) == 0 {
		rows, err = c.doSimpleQuery(ctx, q)
	} else {
		rows, err = c.doExtendedQuery(ctx, q, args)
	}

	if err != nil {
		logger.Logger.DebugError("error in Context Query", err)
		return nil, err
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
		resultSets: []RowsResultSet{
			{
				rowDescriptionBytes: res.RowDescriptionBytes,
				dataRowBytes:        res.DataRowBytes,
			},
		},
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

	resultSets := make([]RowsResultSet, 0, len(res.ResultSets))
	for _, r := range res.ResultSets {
		resultSets = append(resultSets, RowsResultSet{
			rowDescriptionBytes: r.RowDescriptionBytes,
			dataRowBytes:        r.DataRowBytes,
		})
	}

	rows := &Rows{
		resultSets: resultSets,
	}

	return rows, nil
}
