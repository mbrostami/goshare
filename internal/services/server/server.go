package server

import "github.com/mbrostami/goshare/internal/config"

type Service struct {
	cfg *config.Config
}

func NewService(cfg *config.Config) *Service {
	return &Service{
		cfg: cfg,
	}
}
