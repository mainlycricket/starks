package auth

import (
	"bytes"
	"io"

	"github.com/mainlycricket/starks/internal/types"
	byte_utils "github.com/mainlycricket/starks/pkg/byte-utils"
)

func WriteStartupMessage(w io.Writer, connectionInfo *types.ConnectionInfo) error {
	messageBytes, err := getStartupMessageBytes(connectionInfo)
	if err != nil {
		return err
	}

	if _, err := w.Write(messageBytes); err != nil {
		return err
	}

	return nil
}

func getStartupMessageBytes(c *types.ConnectionInfo) ([]byte, error) {
	message := &bytes.Buffer{}

	dataBytes, err := getStartupDataBytes(c)
	if err != nil {
		return nil, err
	}

	if err := byte_utils.WriteInt32Bytes(message, int32(4+len(dataBytes))); err != nil {
		return nil, err
	}

	if _, err := message.Write(dataBytes); err != nil {
		return nil, err
	}

	return message.Bytes(), nil
}

func getStartupDataBytes(c *types.ConnectionInfo) ([]byte, error) {
	dataBuffer := &bytes.Buffer{}

	if err := byte_utils.WriteInt32Bytes(dataBuffer, 196610); err != nil {
		return nil, err
	}

	if err := byte_utils.WriteKeyValueLine(dataBuffer, "user", c.User); err != nil {
		return nil, err
	}

	if err := byte_utils.WriteKeyValueLine(dataBuffer, "database", c.Database); err != nil {
		return nil, err
	}

	for key, value := range c.Options {
		if err := byte_utils.WriteKeyValueLine(dataBuffer, key, value); err != nil {
			return nil, err
		}
	}

	if err := dataBuffer.WriteByte(0); err != nil {
		return nil, err
	}

	return dataBuffer.Bytes(), nil
}
