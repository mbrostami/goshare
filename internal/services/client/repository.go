package client

import "github.com/mbrostami/goshare/internal/models"

type Repository interface {
	// AddServer adds a new server in client db
	AddServer(server *models.Server) error
	GetServer(addr string) (*models.Server, error)
	AddUser(user *models.User) error
	GetUser(username string) (*models.User, error)
}
