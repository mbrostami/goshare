package database

import (
	"fmt"
)

func (r *Repository) AddUser(username, pubKey string) error {
	return r.insertIfNotExist(serverBucketKey, fmt.Sprintf("users-%s", username), pubKey)
}

func (r *Repository) GetUser(username string) (string, error) {
	val, err := r.get(serverBucketKey, fmt.Sprintf("users-%s", username))
	if err != nil {
		return "", err
	}

	return val.(string), nil
}
