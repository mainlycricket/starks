package driver

import (
	"context"
	"database/sql/driver"
	"net"
	"strconv"
	"sync/atomic"

	"github.com/mainlycricket/starks/internal/messages/query"
	"github.com/mainlycricket/starks/pkg/logger"
)

type Conn struct {
	connection       net.Conn
	backendKey       []byte
	backendProcessId int32
	statementCounter atomic.Uint64
	portalCounter    *atomic.Uint64
}

func (c *Conn) Prepare(q string) (driver.Stmt, error) {
	return c.PrepareContext(context.Background(), q)
}

func (c *Conn) PrepareContext(ctx context.Context, q string) (driver.Stmt, error) {
	statementNum := c.statementCounter.Add(1)
	statementName := strconv.FormatUint(statementNum, 10)

	if err := query.PrepareStatement(ctx, c.connection, q, statementName); err != nil {
		return nil, err
	}

	return &Stmt{
		statementName: statementName,
		portalCounter: c.portalCounter,
		c:             c.connection,
	}, nil
}

func (c *Conn) Close() error {
	return nil
}

func (c *Conn) Begin() (driver.Tx, error) {
	return &Tx{}, nil
}

func (c *Conn) Ping(ctx context.Context) error {
	logger.Logger.DebugInfo("Ping invoked")

	if _, err := c.QueryContext(ctx, "SELECT 1;", nil); err != nil {
		logger.Logger.DebugError("error in Simple Ping", "err", err)
		return err
	}

	logger.Logger.DebugInfo("Ping finished")

	return nil
}

type Connector struct {
	conn   *Conn
	driver *Driver
}

func (c *Connector) Connect(ctx context.Context) (driver.Conn, error) {
	return &Conn{}, nil
}

func (c *Connector) Driver() driver.Driver {
	return &Driver{}
}
