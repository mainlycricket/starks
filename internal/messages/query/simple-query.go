package query

import (
	"bytes"
	"context"
	"errors"
	"io"
	"slices"

	"github.com/mainlycricket/starks/internal/messages"
	"github.com/mainlycricket/starks/internal/types"
	byte_utils "github.com/mainlycricket/starks/pkg/byte-utils"
	"github.com/mainlycricket/starks/pkg/operators"
)

func WriteSimpleQuery(ctx context.Context, w io.Writer, query string) error {
	b, err := getSimpleQueryBytes(ctx, query)
	if err != nil {
		return err
	}
	if _, err := w.Write(b); err != nil {
		return err
	}
	return nil
}

func ReadSimpleQueryResp(_ context.Context, r io.Reader) (*types.SimpleQueryRespData, error) {
	state, bytesToRead := messages.MsgInitState, 1
	var messageType byte
	data := &bytes.Buffer{}
	res := &types.SimpleQueryRespData{}

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
			if messageType != 'T' && messageType != 'D' && messageType != 'C' && messageType != 'Z' {
				return nil, errors.New("invalid message type")
			}
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
			data.Reset()
		case messages.MsgDataState:
			bytesToRead -= n
			if bytesToRead == 0 {
				b := slices.Clone(data.Bytes())
				switch messageType {
				case 'T':
					res.RowDescriptionBytes = b
				case 'D':
					res.DataRowBytes = append(res.DataRowBytes, b)
				}
				bytesToRead = operators.Ternary(messageType == 'Z', 0, 1)
				state = operators.Ternary(messageType == 'Z', messages.MsgDoneState, messages.MsgInitState)
				data.Reset()
			}
		}
	}

	return res, nil
}

func getSimpleQueryBytes(_ context.Context, query string) ([]byte, error) {
	queryBytes := byte_utils.GetTerminatedBytes(query)

	msg, err := messages.PrepareFrontendBytes('Q', queryBytes)
	if err != nil {
		return nil, err
	}

	return msg, nil
}
