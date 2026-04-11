package sasl

import (
	"io"

	"github.com/mainlycricket/starks/internal/types"
)

func DoScram256Auth(wr io.ReadWriter, connectionInfo *types.ConnectionInfo) (int32, []byte, error) {
	clientNonce, err := writeClientInitialResponse(wr, connectionInfo.User)
	if err != nil {
		return -1, nil, err
	}

	continueResp, err := readServerAuthContinueMsg(wr)
	if err != nil {
		return -1, nil, err
	}

	expectedSign, err := writeAuthFinalResponse(wr, connectionInfo, clientNonce, continueResp)
	if err != nil {
		return -1, nil, err
	}

	if err := readAuthFinalMessage(wr, expectedSign); err != nil {
		return -1, nil, err
	}

	if err := readAuthOkMessage(wr); err != nil {
		return -1, nil, err
	}

	processId, backendKey, err := readPostAuthMessages(wr)
	if err != nil {
		return -1, nil, err
	}

	return processId, backendKey, nil
}
