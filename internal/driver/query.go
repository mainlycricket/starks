package driver

import (
	"context"
	"database/sql/driver"

	"github.com/mainlycricket/starks/internal/messages"
	"github.com/mainlycricket/starks/internal/messages/query"
	"github.com/mainlycricket/starks/pkg/logger"
)

func (c *Conn) QueryContext(ctx context.Context, q string, args []driver.NamedValue) (driver.Rows, error) {
	logger.Logger.DebugInfo("Context Query invoked")

	c.lock.Lock()
	defer c.lock.Unlock()

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

	rowDesc, err := query.ReadExtendedQueryResp(ctx, c.connection)
	if err != nil {
		return nil, err
	}

	return &Rows{
		c:                   c.connection,
		rowDescriptionBytes: rowDesc,
		lock:                c.lock,
	}, nil
}

func (c *Conn) doSimpleQuery(ctx context.Context, q string) (driver.Rows, error) {
	if err := query.WriteSimpleQuery(ctx, c.connection, q); err != nil {
		return nil, err
	}

	data, err := messages.ReadBackendMessageByType(c.connection, 'T')
	if err != nil {
		return nil, err
	}

	return &Rows{
		rowDescriptionBytes: data,
		c:                   c.connection,
		lock:                c.lock,
	}, nil
}
