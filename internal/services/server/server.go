package server

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) RegisterUser(username, signature, pubKey string) error {
	// TODO verify the signature
	// TODO check if user exists in the system with the same pubkey
	// TODO for now we only support 1 pubkey per username

	// TODO register user into the system with a pair of username -> pubkey and pubkey -> username

	log.Debug().Msgf("registering user: %s", username)

	return s.repo.AddUser(username, pubKey)
}

func (s *Service) GetUser(username string) {
	pubKey, err := s.repo.GetUser(username)
	if err != nil {
		// TODO return err
	}

	// TODO
	fmt.Printf("pubkey %v", pubKey)
}
