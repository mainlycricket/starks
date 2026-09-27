package driver

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/mainlycricket/starks/internal/messages/auth"
	"github.com/mainlycricket/starks/internal/messages/auth/sasl"
	"github.com/mainlycricket/starks/internal/types"
	"github.com/mainlycricket/starks/pkg/logger"
)

func init() {
	sql.Register("starks", Driver{})
}

type Driver struct {
}

func (d Driver) Open(connectionString string) (driver.Conn, error) {
	logger.Logger.DebugInfo("Driver Open invoked")

	connectionInfo, err := d.parseConnectionString(connectionString)
	if err != nil {
		return nil, err
	}

	address := net.JoinHostPort(connectionInfo.Host, fmt.Sprintf("%d", connectionInfo.Port))
	c, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	conn := &Conn{
		connection: c,
	}

	if err := auth.WriteStartupMessage(conn.connection, connectionInfo); err != nil {
		return nil, err
	}

	if _, err = auth.ReadAuthReqMessage(conn.connection); err != nil {
		return nil, err
	}

	processId, backendKey, err := sasl.DoScram256Auth(conn.connection, connectionInfo)
	if err != nil {
		return nil, err
	}

	conn.backendProcessId = processId
	conn.backendKey = backendKey

	logger.Logger.DebugInfo("Driver Open finished")
	return conn, nil
}

func (d *Driver) parseConnectionString(str string) (*types.ConnectionInfo, error) {
	res := &types.ConnectionInfo{
		Options: map[string]string{},
	}

	for part := range strings.SplitSeq(str, ",") {
		partSplit := strings.SplitN(part, "=", 2)
		if len(partSplit) != 2 {
			return nil, ErrInvalidConnectionString
		}
		key, value := partSplit[0], partSplit[1]
		switch key {
		case "host":
			res.Host = value
		case "port":
			parsed, err := strconv.Atoi(value)
			if err != nil {
				return nil, ErrInvalidConnectionString
			}
			res.Port = parsed
		case "user":
			res.User = value
		case "password":
			res.Password = value
		case "database":
			res.Database = value
		default:
			res.Options[key] = value
		}
	}

	return res, nil
}

func (d *Driver) OpenConnector(name string) (driver.Connector, error) {
	return &Connector{}, nil
}
