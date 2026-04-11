package driver

import (
	"context"
	"database/sql/driver"
	"log"
	"net"
)

type Conn struct {
	connection       net.Conn
	backendKey       []byte
	backendProcessId int32
}

func (c *Conn) Prepare(query string) (driver.Stmt, error) {
	return &Stmt{}, nil
}

func (c *Conn) Close() error {
	return nil
}

func (c *Conn) Begin() (driver.Tx, error) {
	return &Tx{}, nil
}

func (c *Conn) Ping(ctx context.Context) error {
	log.Println("Conn Ping invoked")

	if _, err := c.QueryContext(ctx, "SELECT 1;", nil); err != nil {
		log.Println("error in Conn Ping", "err", err)
		return err
	}

	log.Println("Conn Ping finished")

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
