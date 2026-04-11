package sasl

import (
	"bytes"
	"crypto/hmac"
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"strconv"
	"strings"

	"github.com/mainlycricket/starks/internal/messages"
	"github.com/mainlycricket/starks/internal/types"
	byte_utils "github.com/mainlycricket/starks/pkg/byte-utils"
)

func buildClientInitialResponse(username string) ([]byte, string, error) {
	msg, clientNonce := getClientInitialMsg(username)

	data := &bytes.Buffer{}
	if _, err := data.Write(byte_utils.GetTerminatedBytes("SCRAM-SHA-256")); err != nil {
		return nil, "", err
	}

	if err := byte_utils.WriteInt32Bytes(data, int32(len(msg))); err != nil {
		return nil, "", err
	}

	if _, err := data.Write(msg); err != nil {
		return nil, "", err
	}

	msg, err := messages.PrepareFrontendBytes('p', data.Bytes())
	if err != nil {
		return nil, "", err
	}

	return msg, clientNonce, nil
}

func getClientInitialMsg(username string) ([]byte, string) {
	clientNonce := rand.Text()

	bareMessage := getClientInitialBareMsg(username, clientNonce)

	parts := []string{
		"n",
		"",
		string(bareMessage),
	}

	str := strings.Join(parts, ",")
	return []byte(str), clientNonce
}

func buildClientFinalMessage(connectionInfo *types.ConnectionInfo, clientNonce string, saslContinueResp *types.AuthSASLContinueResp) ([]byte, []byte, error) {
	authMessage := getAuthMessage(connectionInfo, clientNonce, saslContinueResp)

	saltedPassword, clientKey, storedKey, err := getPasswordKeys(connectionInfo.Password, saslContinueResp.Salt, saslContinueResp.Iterations)
	if err != nil {
		return nil, nil, err
	}

	proof, err := getClientProof(storedKey, clientKey, authMessage)
	if err != nil {
		return nil, nil, err
	}

	expectedSign, err := getExpectedServerSign(saltedPassword, authMessage)
	if err != nil {
		return nil, nil, err
	}

	first := []byte("c=biws")
	second := fmt.Appendf(nil, "r=%s", saslContinueResp.Nonce)
	third := []byte{'p', '='}
	third = append(third, proof...)

	data := bytes.Join([][]byte{first, second, third}, []byte(","))

	msg, err := messages.PrepareFrontendBytes('p', data)
	if err != nil {
		return nil, nil, err
	}

	return msg, expectedSign, nil
}

func getAuthMessage(connectionInfo *types.ConnectionInfo, clientNonce string, saslContinueResp *types.AuthSASLContinueResp) []byte {
	first := getClientInitialBareMsg(connectionInfo.User, clientNonce)
	second := getServerSASLContinueBareMsg(saslContinueResp)
	third := getClientFinalBareMsgWithoutProof(saslContinueResp.Nonce)

	authMessage := bytes.Join([][]byte{first, second, third}, []byte(","))
	return authMessage
}

func getPasswordKeys(password string, salt []byte, iterations int) (saltedPassword, clientKey []byte, storedKey [32]byte, err error) {
	decodedSalt, err := byte_utils.DecodeBase64(salt, true)
	if err != nil {
		return nil, nil, [32]byte{}, err
	}

	saltedPassword, err = pbkdf2.Key(sha256.New, password, decodedSalt, iterations, 32)
	if err != nil {
		return nil, nil, [32]byte{}, err
	}

	h := hmac.New(sha256.New, saltedPassword)
	if _, err := h.Write([]byte("Client Key")); err != nil {
		return nil, nil, [32]byte{}, err
	}

	clientKey = h.Sum(nil)
	storedKey = sha256.Sum256(clientKey)
	return saltedPassword, clientKey, storedKey, nil
}

func getClientProof(storedKey [32]byte, clientKey, authMesssage []byte) ([]byte, error) {
	clientSignature, err := getClientSignature(authMesssage, storedKey)
	if err != nil {
		return nil, err
	}

	clientProof := byte_utils.XorBytes(clientKey, clientSignature)
	encodedClientProof := byte_utils.EncodeToBase64(clientProof)

	return encodedClientProof, nil
}

func getExpectedServerSign(saltedPassword, authMessage []byte) ([]byte, error) {
	h := hmac.New(sha256.New, saltedPassword)
	if _, err := h.Write([]byte("Server Key")); err != nil {
		return nil, err
	}

	serverKey := h.Sum(nil)

	h2 := hmac.New(sha256.New, serverKey)
	if _, err := h2.Write(authMessage); err != nil {
		return nil, err
	}

	serverSign := h2.Sum(nil)

	base64EncodedSign := byte_utils.EncodeToBase64(serverSign)
	return base64EncodedSign, nil
}

func getClientSignature(authMessage []byte, serverKey [32]byte) ([]byte, error) {
	h := hmac.New(sha256.New, serverKey[:])
	if _, err := h.Write(authMessage); err != nil {
		return nil, err
	}

	clientSignature := h.Sum(nil)

	return clientSignature, nil
}

func getClientInitialBareMsg(username, nonce string) []byte {
	u := []byte{'n', '='}
	u = append(u, []byte(username)...)

	n := []byte{'r', '='}
	n = append(n, []byte(nonce)...)

	return bytes.Join([][]byte{u, n}, []byte{','})
}

func getServerSASLContinueBareMsg(data *types.AuthSASLContinueResp) []byte {
	first := []byte{'r', '='}
	first = append(first, []byte(data.Nonce)...)

	second := []byte{'s', '='}
	second = append(second, data.Salt...)

	third := []byte{'i', '='}
	third = append(third, []byte(strconv.Itoa(data.Iterations))...)
	return bytes.Join([][]byte{first, second, third}, []byte{','})
}

func getClientFinalBareMsgWithoutProof(combinedNonce string) []byte {
	b := [][]byte{
		[]byte("c=biws"),
		fmt.Appendf(nil, "r=%s", combinedNonce),
	}

	return bytes.Join(b, []byte(","))
}
