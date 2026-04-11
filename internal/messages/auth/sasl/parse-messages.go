package sasl

import (
	"bytes"
	"errors"
	"strconv"

	"github.com/mainlycricket/starks/internal/types"
	byte_utils "github.com/mainlycricket/starks/pkg/byte-utils"
)

func parseSASLContinueResp(d []byte) (*types.AuthSASLContinueResp, error) {
	if n, err := byte_utils.ParseInt32Bytes(d[:4]); err != nil || n != 11 {
		return nil, errors.New("invalid value for SASL challenge")
	}

	messageBytes := d[4:]
	res := &types.AuthSASLContinueResp{}

	for part := range bytes.SplitSeq(messageBytes, []byte(",")) {
		lineParts := bytes.SplitN(part, []byte("="), 2)
		if len(lineParts) != 2 {
			return nil, errors.New("invalid data")
		}
		key, value := lineParts[0], lineParts[1]
		switch string(key) {
		case "r":
			res.Nonce = string(value)
		case "s":
			res.Salt = value
		case "i":
			parsed, err := strconv.Atoi(string(value))
			if err != nil {
				return nil, err
			}
			res.Iterations = parsed
		default:
			return nil, errors.New("invalid key")
		}
	}

	return res, nil
}
