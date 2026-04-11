package sasl

import (
	"io"

	"github.com/mainlycricket/starks/internal/types"
)

func writeClientInitialResponse(w io.Writer, username string) (string, error) {
	dataBytes, clientNonce, err := buildClientInitialResponse(username)
	if err != nil {
		return "", err
	}

	if _, err := w.Write(dataBytes); err != nil {
		return "", err
	}

	return clientNonce, nil
}

func writeAuthFinalResponse(w io.Writer, connectionInfo *types.ConnectionInfo, clientNonce string, saslContinueResp *types.AuthSASLContinueResp) ([]byte, error) {
	data, expectedSign, err := buildClientFinalMessage(connectionInfo, clientNonce, saslContinueResp)
	if err != nil {
		return expectedSign, err
	}

	if _, err := w.Write(data); err != nil {
		return expectedSign, err
	}

	return expectedSign, nil
}
