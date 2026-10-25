package driver

import (
	"context"
	"database/sql/driver"
	"net"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/mainlycricket/starks/internal/messages/query"
)

type Stmt struct {
	c             net.Conn    // points to the original TCP connection to the database
	lock          sync.Locker // points to the original lock
	statementName string
	portalCounter *atomic.Uint64
}

func (s *Stmt) Close() error {
	return nil
}

func (s *Stmt) NumInput() int {
	return -1
}

func (s *Stmt) Exec(args []driver.Value) (driver.Result, error) {
	return s.ExecContext(context.Background(), args)
}

func (s *Stmt) Query(args []driver.Value) (driver.Rows, error) {
	return s.QueryContext(context.Background(), args)
}

func (s *Stmt) ExecContext(ctx context.Context, args []driver.Value) (driver.Result, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	portalNumber := s.portalCounter.Add(1)
	portalNumberStr := strconv.FormatUint(portalNumber, 10)

	if err := query.WriteExecStatement(ctx, s.c, s.statementName, portalNumberStr, args); err != nil {
		return nil, err
	}

	resp, err := query.ReadExecStmtResp(ctx, s.c)
	if err != nil {
		return nil, err
	}

	return newResult(resp.CommandCompleteData), nil
}

func (s *Stmt) QueryContext(ctx context.Context, args []driver.Value) (driver.Rows, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	portalNumber := s.portalCounter.Add(1)
	portalNumberStr := strconv.FormatUint(portalNumber, 10)

	if err := query.WriteExecStatement(ctx, s.c, s.statementName, portalNumberStr, args); err != nil {
		return nil, err
	}

	rowDesc, err := query.ReadStmtQueryResp(ctx, s.c)
	if err != nil {
		return nil, err
	}

	return &Rows{
		rowDescriptionBytes: rowDesc,
		c:                   s.c,
		lock:                s.lock,
	}, nil
}
