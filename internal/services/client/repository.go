package client

type Repository interface {
	// AddServer adds a new server in client db
	AddServer(addr string) error
}
