package types

type AuthReqInfo struct {
	AuthType int32
	Methods  []string
}

type AuthSASLContinueResp struct {
	Nonce      string
	Salt       []byte // base64 encoded
	Iterations int
}
