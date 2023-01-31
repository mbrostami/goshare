package server

import "github.com/mbrostami/goshare/internal/models"

type Repository interface {
	AddUserToServer(user *models.User) error
	GetUserFromServer(username string) (*models.User, error)
}
