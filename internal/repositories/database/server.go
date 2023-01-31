package database

import (
	"fmt"

	"github.com/mbrostami/goshare/internal/models"
)

func (r *Repository) AddUserToServer(user *models.User) error {
	return r.insertIfNotExist(serverBucketKey, fmt.Sprintf("users-%s", user.Username), user)
}

func (r *Repository) GetUserFromServer(username string) (*models.User, error) {
	var usr models.User
	err := r.get(serverBucketKey, fmt.Sprintf("users-%s", username), &usr)
	if err != nil {
		return nil, err
	}

	return &usr, nil
}
