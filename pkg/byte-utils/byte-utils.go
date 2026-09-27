package byte_utils

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"io"
	"strconv"
)

func WriteInt32Bytes(w io.Writer, value int32) error {
	valueBytes, err := GetIntBytes(value)
	if err != nil {
		return err
	}

	if _, err := w.Write(valueBytes); err != nil {
		return err
	}

	return nil
}

func WriteKeyValueLine(w io.Writer, key, value string) error {
	lineBytes := GetKeyValueLine(key, value)
	if _, err := w.Write(lineBytes); err != nil {
		return err
	}
	return nil
}

func GetIntBytes[T int16 | int32](value T) ([]byte, error) {
	buf := &bytes.Buffer{}
	if err := binary.Write(buf, binary.BigEndian, value); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func GetKeyValueLine(key, value string) []byte {
	n := len(key) + len(value) + 2

	b := make([]byte, 0, n)

	b = append(b, []byte(key)...)
	b = append(b, 0)
	b = append(b, []byte(value)...)
	b = append(b, 0)

	return b
}

func ParseInt32Bytes(d []byte) (int32, error) {
	var value int32
	_, err := binary.Decode(d, binary.BigEndian, &value)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func ParseInt16Bytes(d []byte) (int16, error) {
	var value int16
	_, err := binary.Decode(d, binary.BigEndian, &value)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func GetTerminatedBytes(s string) []byte {
	b := []byte(s)
	b = append(b, 0)
	return b
}

func XorBytes(a, b []byte) []byte {
	n := min(len(a), len(b))
	res := make([]byte, n)

	for i := range n {
		res[i] = a[i] ^ b[i]
	}
	return res
}

func EncodeToBase64(d []byte) []byte {
	n := base64.StdEncoding.EncodedLen(len(d))
	res := make([]byte, n)
	base64.StdEncoding.Encode(res, d)
	return res
}

func DecodeBase64(d []byte, trimZeroBytes bool) ([]byte, error) {
	n := base64.StdEncoding.DecodedLen(len(d))
	res := make([]byte, n)

	if _, err := base64.StdEncoding.Decode(res, d); err != nil {
		return nil, err
	}

	if trimZeroBytes {
		res = bytes.TrimRight(res, string([]byte{0}))
	}

	return res, nil
}

// supported for int, float, bool, string
func GetPostgresTextBytes(v any) ([]byte, error) {
	if v == nil {
		return []byte{}, nil
	}

	switch v := v.(type) {
	case int:
		return []byte(strconv.Itoa(v)), nil
	case int8:
		return []byte(strconv.FormatInt(int64(v), 10)), nil
	case int16:
		return []byte(strconv.FormatInt(int64(v), 10)), nil
	case int32:
		return []byte(strconv.FormatInt(int64(v), 10)), nil
	case int64:
		return []byte(strconv.FormatInt(v, 10)), nil
	case float64:
		return []byte(strconv.FormatFloat(v, 'g', -1, 64)), nil
	case float32:
		return []byte(strconv.FormatFloat(float64(v), 'g', -1, 32)), nil
	case bool:
		return []byte(strconv.FormatBool(v)), nil
	case string:
		return []byte(v), nil
	default:
		return nil, errors.New("unsupported type")
	}
}
