package auth

import (
	"bytes"
	"errors"
	"io"
	"slices"

	"github.com/mainlycricket/starks/internal/messages"
	"github.com/mainlycricket/starks/internal/types"
	byte_utils "github.com/mainlycricket/starks/pkg/byte-utils"
)

const (
	AuthReqTypeSASL int32 = 10
)

func ReadAuthReqMessage(r io.Reader) (*types.AuthReqInfo, error) {
	data, err := messages.ReadBackendBytes(r, 'R')
	if err != nil {
		return nil, err
	}

	authReqInfo, err := parseAuthReqInfo(data)
	if err != nil {
		return nil, err
	}

	return authReqInfo, nil
}

func parseAuthReqInfo(d []byte) (*types.AuthReqInfo, error) {
	authType, err := byte_utils.ParseInt32Bytes(d[:4])
	if err != nil {
		return nil, err
	}

	if authType != AuthReqTypeSASL {
		return nil, errors.New("auth mechanism currently not supported")
	}

	res := &types.AuthReqInfo{
		AuthType: AuthReqTypeSASL,
	}

	for method := range bytes.SplitSeq(d[4:len(d)-1], []byte{0}) {
		if len(method) == 0 {
			continue
		}
		res.Methods = append(res.Methods, string(method))
	}

	if !slices.Contains(res.Methods, "SCRAM-SHA-256") {
		return nil, errors.New("auth method not supported")
	}

	return res, nil
}
