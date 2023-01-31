package database

import (
	"encoding/json"
	"errors"

	"github.com/boltdb/bolt"
)

var clientBucketKey = []byte("client")
var serverBucketKey = []byte("server")

type Repository struct {
	db           *bolt.DB
	serverBucket *bolt.Bucket
	clientBucket *bolt.Bucket
}

func NewRepository(db *bolt.DB) (*Repository, error) {
	var repo Repository

	err := db.Update(func(tx *bolt.Tx) error {
		s, err := tx.CreateBucketIfNotExists(serverBucketKey)
		if err != nil {
			return err
		}
		repo.serverBucket = s

		c, err := tx.CreateBucketIfNotExists(clientBucketKey)
		if err != nil {
			return err
		}
		repo.clientBucket = c
		return nil
	})

	if err != nil {
		return nil, err
	}

	repo.db = db

	return &repo, nil
}

func (r *Repository) insertAndOverride(bucketKey []byte, key string, value interface{}) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		jsonString, err := json.Marshal(value)
		if err != nil {
			return err
		}
		return tx.Bucket(bucketKey).Put([]byte(key), jsonString)
	})
}

func (r *Repository) insertIfNotExist(bucketKey []byte, key string, value interface{}) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		if v := tx.Bucket(bucketKey).Get([]byte(key)); v != nil {
			return errors.New("key already exist")
		}
		jsonString, err := json.Marshal(value)
		if err != nil {
			return err
		}
		return tx.Bucket(bucketKey).Put([]byte(key), jsonString)
	})
}

func (r *Repository) get(bucketKey []byte, key string) (interface{}, error) {
	var value interface{}
	err := r.db.View(func(tx *bolt.Tx) error {
		if v := tx.Bucket(bucketKey).Get([]byte(key)); v != nil {
			err := json.Unmarshal(v, &value)
			if err != nil {
				return err
			}
		}
		return errors.New("key doesn't exist")
	})
	if err != nil {
		return nil, err
	}
	return value, nil
}
