package client

import (
	"context"

	"github.com/mbrostami/goshare/api/grpc"
	"github.com/rs/zerolog/log"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) AddServer(addr string) error {
	return s.repo.AddServer(addr)
}

func (s *Service) Register(username, addr, pubKey string) error {
	c, err := grpc.NewClient(addr)
	if err != nil {
		return err
	}

	log.Debug().Msgf("registering user: %s", username)

	return c.Register(context.Background(), username, pubKey)
}
