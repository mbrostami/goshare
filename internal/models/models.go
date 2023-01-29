package models

type Server struct {
	Address string
	Auth    string
}

type User struct {
	Username string
	PubKey   []byte
}

type File struct {
	Metadata string
	Name     string
	Size     int
	Checksum string
	Chunks   [][]byte
}
