package types

type ConnectionInfo struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	Options  map[string]string
}
