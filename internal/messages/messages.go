package messages

import (
	"bytes"
	"errors"
	"io"
	"slices"
	"strings"

	byte_utils "github.com/mainlycricket/starks/pkg/byte-utils"
	"github.com/mainlycricket/starks/pkg/operators"
)

const (
	MsgInitState = iota
	MsgLengthState
	MsgDataState
	MsgDoneState
)

func PrepareFrontendBytes(msgType byte, data []byte) ([]byte, error) {
	resp := &bytes.Buffer{}

	if err := resp.WriteByte(msgType); err != nil {
		return nil, err
	}

	if err := byte_utils.WriteInt32Bytes(resp, int32(len(data)+4)); err != nil {
		return nil, err
	}

	if _, err := resp.Write(data); err != nil {
		return nil, err
	}

	return resp.Bytes(), nil
}

func ReadBackendBytes(r io.Reader, msgType byte) ([]byte, error) {
	state, bytesToRead := MsgInitState, 1
	data := &bytes.Buffer{}
	var isError bool

	for state != MsgDoneState {
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
		case MsgInitState:
			if b[0] == 'E' {
				isError = true
			} else if b[0] != msgType {
				return nil, errors.New("invalid first byte")
			}
			state, bytesToRead = MsgLengthState, 4
			data.Reset()
		case MsgLengthState:
			if data.Len() != 4 {
				continue
			}
			messageSize, err := byte_utils.ParseInt32Bytes(data.Bytes())
			if err != nil {
				return nil, err
			}
			state, bytesToRead = MsgDataState, int(messageSize)-4
			data.Reset()
		case MsgDataState:
			bytesToRead -= n
			state = operators.Ternary(bytesToRead == 0, MsgDoneState, MsgDataState)
		}
	}

	if isError {
		err := prepareError(data.Bytes())
		return nil, err
	}

	return slices.Clone(data.Bytes()), nil
}

// TODO: analyze / classify / structure / prettify error
func prepareError(errBytes []byte) error {
	if len(errBytes) == 1 || errBytes[0] == '0' {
		return errors.New("something went wrong")
	}

	var parts []string

	for b := range bytes.SplitSeq(errBytes[1:], []byte{0}) {
		if len(b) > 0 {
			parts = append(parts, string(b))
		}
	}

	msg := strings.Join(parts, " ")

	return errors.New(msg)
}
