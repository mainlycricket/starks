package query

import (
	"context"
	"database/sql/driver"
	"io"

	"github.com/mainlycricket/starks/internal/messages"
	"github.com/mainlycricket/starks/internal/types"
	"github.com/mainlycricket/starks/pkg/logger"
)

func PrepareStatement(ctx context.Context, wr io.ReadWriter, q, statementName string) error {
	logger.Logger.DebugInfo("query PrepareStatement invoked")

	if err := writeParseCommand(ctx, wr, q, statementName); err != nil {
		logger.Logger.DebugInfo("error in writeParseCommand in prepare stmt", err)
		return err
	}

	if err := writeSyncCommand(ctx, wr); err != nil {
		logger.Logger.DebugInfo("error in writeSyncCommand in prepare stmt", err)
		return err
	}

	if err := readParseResponse(ctx, wr); err != nil {
		logger.Logger.DebugInfo("error in readParseResponse in prepare stmt", err)
		return err
	}

	if err := readSyncResponse(ctx, wr); err != nil {
		logger.Logger.DebugInfo("error in readSyncResponse in prepare stmt", err)
		return err
	}

	logger.Logger.DebugInfo("query PrepareStatement finished")

	return nil
}

func WriteExecStatement(ctx context.Context, wr io.ReadWriter, statementName, portalName string, args []driver.Value) error {
	logger.Logger.DebugInfo("query WriteExecStatement invoked")

	argValues := make([]driver.NamedValue, 0, len(args))
	for _, arg := range args {
		argValues = append(argValues, driver.NamedValue{Value: arg})
	}

	if err := writeBindCommand(ctx, wr, statementName, portalName, argValues); err != nil {
		logger.Logger.DebugInfo("error in writeBindCommand in exec stmt", err)
		return err
	}

	if err := writeDescribeCommand(ctx, wr, statementName); err != nil {
		logger.Logger.DebugInfo("error in writeDescribeCommand in exec stmt", err)
		return err
	}

	if err := writeExecuteCommand(ctx, wr, portalName); err != nil {
		logger.Logger.DebugInfo("error in writeExecuteCommand in exec stmt", err)
		return err
	}

	if err := writeSyncCommand(ctx, wr); err != nil {
		logger.Logger.DebugInfo("error in writeSyncCommand in exec stmt", err)
		return err
	}

	logger.Logger.DebugInfo("query WriteExecStatement finished")

	return nil
}

func ReadExecStmtResp(ctx context.Context, r io.ReadWriter) (*types.ExtenededQueryRespData, error) {
	if err := readBindResponse(ctx, r); err != nil {
		logger.Logger.DebugInfo("error in readBindResponse in exec stmt", err)
		return nil, err
	}

	return readExecResponse(ctx, r)
}

func ReadStmtQueryResp(ctx context.Context, r io.ReadWriter) ([]byte, error) {
	if err := readBindResponse(ctx, r); err != nil {
		return nil, err
	}

	if _, err := messages.ReadBackendMessageByType(r, 't'); err != nil {
		return nil, err
	}

	data, err := messages.ReadBackendMessageByType(r, 'T')
	if err != nil {
		return nil, err
	}

	return data, nil
}
