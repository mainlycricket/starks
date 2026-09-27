package types

type SimpleQueryRespData struct {
	RowDescriptionBytes []byte
	DataRowBytes        [][]byte
}

type ExtenededQueryRespData struct {
	RowDescriptionBytes  []byte
	DataRowBytes         [][]byte
	ParameterDescription []byte
}
