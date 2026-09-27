package driver

import (
	"context"
	"database/sql/driver"

	"github.com/mainlycricket/starks/internal/messages/query"
	"github.com/mainlycricket/starks/pkg/logger"
)

func (c *Conn) ExecContext(ctx context.Context, q string, args []driver.NamedValue) (driver.Result, error) {
	logger.Logger.DebugInfo("ExecContext invoked")

	if err := query.WriteExtendedQuery(ctx, c.connection, q, args); err != nil {
		return nil, err
	}

	res, err := query.ReadExtendedQueryResp(ctx, c.connection)
	if err != nil {
		return nil, err
	}

	logger.Logger.DebugInfo("ExecContext finished")
	return newResult(res.CommandCompleteData), nil
}
