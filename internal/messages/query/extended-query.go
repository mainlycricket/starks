package query

import (
	"bytes"
	"context"
	"database/sql/driver"
	"io"
	"slices"

	"github.com/mainlycricket/starks/internal/messages"
	"github.com/mainlycricket/starks/internal/types"
	byte_utils "github.com/mainlycricket/starks/pkg/byte-utils"
	"github.com/mainlycricket/starks/pkg/logger"
	"github.com/mainlycricket/starks/pkg/operators"
)

func WriteExtendedQuery(ctx context.Context, w io.Writer, query string, args []driver.NamedValue) error {
	if err := writeParseCommand(ctx, w, query); err != nil {
		return err
	}

	if err := writeBindCommand(ctx, w, args); err != nil {
		return err
	}

	if err := writeDescribeCommand(ctx, w); err != nil {
		return err
	}

	if err := writeExecuteCommand(ctx, w); err != nil {
		return err
	}

	if err := writeSyncCommand(ctx, w); err != nil {
		return err
	}

	return nil
}

func ReadExtendedQueryResp(ctx context.Context, r io.Reader) (*types.ExtenededQueryRespData, error) {
	if err := readParseResponse(ctx, r); err != nil {
		return nil, err
	}

	if err := readBindResponse(ctx, r); err != nil {
		return nil, err
	}

	return readDataRows(ctx, r)
}

func writeParseCommand(_ context.Context, w io.Writer, query string) error {
	logger.Logger.DebugInfo("writeParseCommand invoked")

	stmtName := byte_utils.GetTerminatedBytes("")
	queryBytes := byte_utils.GetTerminatedBytes(query)
	prmtrTypesBytes, err := byte_utils.GetIntBytes(int16(0))
	if err != nil {
		return err
	}

	data := make([]byte, 0, len(stmtName)+len(queryBytes)+len(prmtrTypesBytes))
	data = append(data, stmtName...)
	data = append(data, queryBytes...)
	data = append(data, prmtrTypesBytes...)

	msgBytes, err := messages.PrepareFrontendBytes('P', data)
	if err != nil {
		return err
	}

	logger.Logger.DebugInfo("message bytes in writeParseCommand", msgBytes)

	if _, err := w.Write(msgBytes); err != nil {
		return err
	}

	logger.Logger.DebugInfo("writeParseCommand finished")

	return nil
}

func writeBindCommand(_ context.Context, w io.Writer, args []driver.NamedValue) error {
	logger.Logger.DebugInfo("writeBindCommand invoked")

	destinationPortalBytes := byte_utils.GetTerminatedBytes("")
	stmtNameBytes := byte_utils.GetTerminatedBytes("")
	prmtrFmtCodeBytes, _ := byte_utils.GetIntBytes(int16(0))

	prmtrCountBytes, _ := byte_utils.GetIntBytes(int16(len(args)))
	var prmtrBytes []byte
	for _, arg := range args {
		argValueBytes, err := byte_utils.GetPostgresTextBytes(arg.Value)
		if err != nil {
			return err
		}

		argLen := operators.Ternary(arg.Value == nil, -1, len(argValueBytes))
		argLenBytes, _ := byte_utils.GetIntBytes(int32(argLen))

		prmtrBytes = append(prmtrBytes, argLenBytes...)
		prmtrBytes = append(prmtrBytes, argValueBytes...)
	}

	resFmtCodeBytes, _ := byte_utils.GetIntBytes(int16(0))

	dataBytes := bytes.Join([][]byte{
		destinationPortalBytes, stmtNameBytes, prmtrFmtCodeBytes, prmtrCountBytes, prmtrBytes, resFmtCodeBytes,
	}, nil)

	msgBytes, err := messages.PrepareFrontendBytes('B', dataBytes)
	if err != nil {
		return err
	}

	logger.Logger.DebugInfo("message bytes in writeBindCommand", msgBytes)
	if _, err := w.Write(msgBytes); err != nil {
		return err
	}

	logger.Logger.DebugInfo("writeBindCommand finished")

	return nil
}

func writeDescribeCommand(_ context.Context, w io.Writer) error {
	logger.Logger.DebugInfo("writeDescribeCommand invoked")

	typeBytes := []byte{'S'}
	statementName := byte_utils.GetTerminatedBytes("")
	payloadBytes := bytes.Join([][]byte{typeBytes, statementName}, nil)

	msgBytes, err := messages.PrepareFrontendBytes('D', payloadBytes)
	if err != nil {
		return err
	}

	logger.Logger.DebugInfo("message bytes in writeDescribeCommand", msgBytes)

	if _, err := w.Write(msgBytes); err != nil {
		return err
	}

	logger.Logger.DebugInfo("writeDescribeCommand finished")

	return nil
}

func writeExecuteCommand(_ context.Context, w io.Writer) error {
	logger.Logger.DebugInfo("writeExecuteCommand invoked")

	portalBytes := byte_utils.GetTerminatedBytes("")
	rowCountBytes, _ := byte_utils.GetIntBytes(int32(0))
	data := bytes.Join([][]byte{portalBytes, rowCountBytes}, nil)
	msgBytes, err := messages.PrepareFrontendBytes('E', data)
	if err != nil {
		return err
	}

	logger.Logger.DebugInfo("message bytes in writeExecuteCommand", msgBytes)
	if _, err := w.Write(msgBytes); err != nil {
		return err
	}

	logger.Logger.DebugInfo("writeExecuteCommand finished")

	return err
}

func writeSyncCommand(_ context.Context, w io.Writer) error {
	logger.Logger.DebugInfo("writeSyncCommand invoked")

	msgBytes, err := messages.PrepareFrontendBytes('S', nil)
	if err != nil {
		return err
	}

	logger.Logger.DebugInfo("message bytes in writeSyncCommand", msgBytes)

	if _, err := w.Write(msgBytes); err != nil {
		return err
	}

	logger.Logger.DebugInfo("writeSyncCommand finished")

	return nil
}

func readParseResponse(_ context.Context, r io.Reader) error {
	logger.Logger.DebugInfo("readParseResponse invoked")

	if _, err := messages.ReadBackendBytes(r, '1'); err != nil {
		logger.Logger.DebugError("while readParseResponse", err)
		return err
	}

	logger.Logger.DebugInfo("readParseResponse finished")
	return nil
}

func readBindResponse(_ context.Context, r io.Reader) error {
	logger.Logger.DebugInfo("readBindResponse invoked")

	if _, err := messages.ReadBackendBytes(r, '2'); err != nil {
		logger.Logger.DebugError("while readBindResponse", err)
		return err
	}

	logger.Logger.DebugInfo("readBindResponse finished")
	return nil
}

func readDataRows(_ context.Context, r io.Reader) (*types.ExtenededQueryRespData, error) {
	state, bytesToRead := messages.MsgInitState, 1
	var messageType byte
	data := &bytes.Buffer{}
	res := &types.ExtenededQueryRespData{}

	for state != messages.MsgDoneState {
		b := make([]byte, bytesToRead)
		n, err := r.Read(b)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		b = b[:n]
		if _, err := data.Write(b); err != nil {
			return nil, err
		}

		switch state {
		case messages.MsgInitState:
			messageType = b[0]
			logger.Logger.DebugInfo("received message type", string(messageType))
			state, bytesToRead = messages.MsgLengthState, 4
			data.Reset()
		case messages.MsgLengthState:
			if data.Len() != 4 {
				continue
			}
			messageSize, err := byte_utils.ParseInt32Bytes(data.Bytes())
			if err != nil {
				return nil, err
			}
			state, bytesToRead = messages.MsgDataState, int(messageSize)-4
			logger.Logger.DebugInfo("message payload length", bytesToRead)
			data.Reset()
		case messages.MsgDataState:
			bytesToRead -= n
			if bytesToRead == 0 {
				payload := slices.Clone(data.Bytes())

				logger.Logger.DebugInfo("message payload", payload)

				switch messageType {
				case 'T':
					res.RowDescriptionBytes = payload
				case 'D':
					res.DataRowBytes = append(res.DataRowBytes, payload)
				case 'C':
					res.CommandCompleteData = payload
				case 't':
					res.ParameterDescription = payload
				}

				state = operators.Ternary(messageType == 'Z', messages.MsgDoneState, messages.MsgInitState)
				bytesToRead = operators.Ternary(state == messages.MsgDoneState, 0, 1)
				data.Reset()
			}
		}
	}

	return res, nil
}
