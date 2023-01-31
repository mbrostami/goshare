package client

import (
	"context"
	"fmt"

	"github.com/mbrostami/goshare/api/grpc"
	"github.com/mbrostami/goshare/internal/models"
	"github.com/mbrostami/goshare/pkg/crypto"
	"github.com/rs/zerolog/log"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) AddServer(addr string) error {
	return s.repo.AddServer(&models.Server{
		Address: addr,
		Auth:    "",
	})
}

// Register adds user to the local db and registers in remote server
// if successful the server will be added to local db as well
func (s *Service) Register(username, addr string) (*models.User, error) {
	var user *models.User
	var err error
	user, err = s.repo.GetUser(username)
	if err != nil {
		pub, priv, err := crypto.GeneratePrivatePubKey()
		if err != nil {
			return nil, err
		}

		user = &models.User{
			Username: username,
			PubKey:   pub,
			PrivKey:  priv,
		}
		log.Debug().Msgf("creating user: %s", user.Username)

		if err := s.repo.AddUser(user); err != nil {
			return nil, err
		}
	}

	c, err := grpc.NewClient(addr)
	if err != nil {
		return nil, fmt.Errorf("connecting to grpc server: %v", err)
	}

	log.Debug().Msgf("registering user: %s", user.Username)

	err = c.Register(context.Background(), username, user.PubKey, crypto.Sign(user.PrivKey, username))
	if err != nil {
		return nil, err
	}

	log.Debug().Msgf("adding server: %s", addr)

	return user, s.AddServer(addr)
}
