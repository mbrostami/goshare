package database

import (
	"fmt"

	"github.com/mbrostami/goshare/internal/models"
)

func (r *Repository) AddServer(server *models.Server) error {
	return r.insertAndOverride(clientBucketKey, fmt.Sprintf("server-%s", server.Address), server)
}

func (r *Repository) GetServer(addr string) (*models.Server, error) {
	var server models.Server
	err := r.get(clientBucketKey, fmt.Sprintf("server-%s", addr), &server)
	if err != nil {
		return nil, err
	}

	return &server, nil
}

func (r *Repository) AddUser(user *models.User) error {
	return r.insertIfNotExist(clientBucketKey, fmt.Sprintf("users-%s", user.Username), user)
}

func (r *Repository) GetUser(username string) (*models.User, error) {
	var usr models.User
	err := r.get(clientBucketKey, fmt.Sprintf("users-%s", username), &usr)
	if err != nil {
		return nil, err
	}

	return &usr, nil
}
