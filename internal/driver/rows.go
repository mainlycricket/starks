package driver

import (
	"database/sql/driver"
	"errors"
	"io"
	"net"
	"sync"

	"github.com/mainlycricket/starks/internal/messages"
	byte_utils "github.com/mainlycricket/starks/pkg/byte-utils"
	"github.com/mainlycricket/starks/pkg/logger"
	"github.com/mainlycricket/starks/pkg/operators"
)

const (
	rowStateNextResSet = iota
	rowStateEnd
)

type Rows struct {
	c    net.Conn    // points to the original TCP connection
	lock sync.Locker // points to the original lock
	done bool

	rowDescriptionBytes     []byte
	nextRowDescriptionBytes []byte
}

type RowsResultSet struct {
	rowDescriptionBytes []byte
	dataRowBytes        [][]byte
	rowIdx              int
}

func (r *Rows) Columns() []string {
	count, _ := byte_utils.ParseInt16Bytes(r.rowDescriptionBytes[:2])
	columns := make([]string, 0, count)

	var columnName []byte
	for idx := 2; idx < len(r.rowDescriptionBytes); {
		b := r.rowDescriptionBytes[idx]
		if b == 0 {
			columns = append(columns, string(columnName))
			idx += 19
			columnName = nil
		} else {
			columnName = append(columnName, b)
			idx++
		}
	}

	return columns
}

func (r *Rows) Close() error {
	logger.Logger.DebugInfo("Rows Close invoked")

	if r.done {
		logger.Logger.DebugInfo("exiting with r done in Rows Close")
		return nil
	}

	if err := r.skipTo(rowStateEnd); err != nil {
		logger.Logger.DebugError("error in skipTo in rows Close", err)
		return err
	}

	r.done = true
	logger.Logger.DebugInfo("Rows Closed finished")
	return nil
}

func (r *Rows) HasNextResultSet() bool {
	logger.Logger.DebugInfo("HasNextResultSet invoked")

	if len(r.nextRowDescriptionBytes) > 0 {
		return true
	}

	if err := r.skipTo(rowStateNextResSet); err != nil {
		logger.Logger.DebugError("error in HasNextResultSet", err)
		return false
	}

	logger.Logger.DebugInfo("HasNextResultSet finished")
	return true
}

func (r *Rows) NextResultSet() error {
	logger.Logger.DebugInfo("NextResultSet invoked")

	if len(r.nextRowDescriptionBytes) > 0 {
		r.rowDescriptionBytes = r.nextRowDescriptionBytes
		r.nextRowDescriptionBytes = nil
		return nil
	}

	if err := r.skipTo(rowStateNextResSet); err != nil {
		logger.Logger.DebugError("error in NextResultSet", err)
		return err
	}

	logger.Logger.DebugInfo("NextResultSet finished")
	return nil
}

func (r *Rows) Next(dest []driver.Value) error {
	logger.Logger.DebugInfo("Next invoked")

	messageType, data, err := r.readBackendMessage()
	if err != nil {
		logger.Logger.DebugError("error in Next", err)
		return err
	}

	switch messageType {
	case 'Z':
		logger.Logger.DebugInfo("Next finished with ReadyForQuery")
		r.done = true
		return io.EOF
	case 'C':
		logger.Logger.DebugInfo("Next finished with Command Complete", string(data))
		return io.EOF
	case 'D':
		var bytesCounter int

		count, _ := byte_utils.ParseInt16Bytes(data[:2])
		bytesCounter += 2

		for i := range count {
			dataLen, _ := byte_utils.ParseInt32Bytes(data[bytesCounter : bytesCounter+4])
			bytesCounter += 4
			if dataLen > 0 {
				dest[i] = data[bytesCounter : bytesCounter+int(dataLen)]
			} else {
				dest[i] = nil
			}
			bytesCounter += operators.Ternary(dataLen > 0, int(dataLen), 0)
		}
		logger.Logger.DebugInfo("Next finished with DataRow")
		return nil
	default:
		logger.Logger.DebugError("Next finished with unexpected type")
		return errors.New("unexpected message type")
	}
}

func (r *Rows) skipTo(desiredState int) error {
	for {
		messageType, data, err := r.readBackendMessage()
		if err != nil {
			return err
		}
		switch messageType {
		case 'E':
			return errors.New(string(data))
		case 'Z':
			r.done = true
			if desiredState == rowStateEnd {
				return nil
			}
			return io.EOF
		case 'C', 'D', 't':
			continue
		case 'T':
			if desiredState == rowStateNextResSet {
				r.nextRowDescriptionBytes = data
				return nil
			}
			continue
		default:
			return errors.New("unexpected message type")
		}
	}
}

func (r *Rows) readBackendMessage() (byte, []byte, error) {
	if r.done {
		return 0, nil, io.EOF
	}

	r.lock.Lock()
	defer r.lock.Unlock()

	messageType, data, err := messages.ReadBackendMessage(r.c)
	if err != nil {
		return 0, nil, err
	}

	if messageType == 'E' {
		return 0, nil, errors.New(string(data))
	}

	return messageType, data, nil
}
