package server

type Repository interface {
	AddUser(username, pubKey string) error
	GetUser(username string) (string, error)
}
