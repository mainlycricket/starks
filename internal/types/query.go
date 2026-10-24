package types

type SimpleQueryRespData struct {
	ResultSets []SimpleQueryResultSet
}

type SimpleQueryResultSet struct {
	RowDescriptionBytes []byte
	DataRowBytes        [][]byte
}

type ExtenededQueryRespData struct {
	Error                []byte
	RowDescriptionBytes  []byte
	DataRowBytes         [][]byte
	ParameterDescription []byte
	CommandCompleteData  []byte
}
