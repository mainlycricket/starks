package sasl

import (
	"bytes"
	"errors"
	"io"

	"github.com/mainlycricket/starks/internal/messages"
	"github.com/mainlycricket/starks/internal/types"
	byte_utils "github.com/mainlycricket/starks/pkg/byte-utils"
)

func readServerAuthContinueMsg(r io.Reader) (*types.AuthSASLContinueResp, error) {
	data, err := messages.ReadBackendMessageByType(r, 'R')

	resp, err := parseSASLContinueResp(data)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func readAuthFinalMessage(r io.Reader, expectedSign []byte) error {
	data, err := messages.ReadBackendMessageByType(r, 'R')
	if err != nil {
		return err
	}

	if n, err := byte_utils.ParseInt32Bytes(data[:4]); err != nil || n != 12 {
		return errors.New("auth not complete")
	}

	if !bytes.Equal(expectedSign, data[6:]) {
		return errors.New("unexpected server sign")
	}

	return nil
}

func readAuthOkMessage(r io.Reader) error {
	data, err := messages.ReadBackendMessageByType(r, 'R')
	if err != nil {
		return err
	}

	if n, err := byte_utils.ParseInt32Bytes(data); err != nil || n != 0 {
		return errors.New("auth not ok")
	}

	return nil
}

func readPostAuthMessages(r io.Reader) (int32, []byte, error) {
	state, bytesToRead := messages.MsgInitState, 1
	data := &bytes.Buffer{}
	var (
		messageType byte
		processId   int32
		backendKey  []byte
	)

	for state != messages.MsgDoneState {
		b := make([]byte, bytesToRead)
		n, err := r.Read(b)
		if err == io.EOF {
			break
		}
		if err != nil {
			return -1, nil, err
		}
		b = b[:n]
		if _, err := data.Write(b); err != nil {
			return -1, nil, err
		}

		switch state {
		case messages.MsgInitState:
			messageType = b[0]
			if messageType != 'S' && messageType != 'K' && messageType != 'Z' {
				return -1, nil, errors.New("invalid first byte")
			}
			state, bytesToRead = messages.MsgLengthState, 4
			data.Reset()
		case messages.MsgLengthState:
			if data.Len() != 4 {
				continue
			}
			messageSize, err := byte_utils.ParseInt32Bytes(data.Bytes())
			if err != nil {
				return -1, nil, err
			}
			state, bytesToRead = messages.MsgDataState, int(messageSize)-4
			data.Reset()
		case messages.MsgDataState:
			bytesToRead -= n
			if bytesToRead == 0 {
				switch messageType {
				case 'S':
					bytesToRead, state = 1, messages.MsgInitState
				case 'K':
					d := data.Bytes()
					parsed, err := byte_utils.ParseInt32Bytes(d[:4])
					if err != nil {
						return -1, nil, err
					}
					processId = parsed
					backendKey = d[4:]
					bytesToRead, state = 1, messages.MsgInitState
				case 'Z':
					state = messages.MsgDoneState
				}
				data.Reset()
			}
		}
	}

	return processId, backendKey, nil
}
