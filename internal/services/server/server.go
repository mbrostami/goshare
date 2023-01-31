package server

import (
	"bytes"
	"errors"

	"github.com/mbrostami/goshare/pkg/crypto"

	"github.com/mbrostami/goshare/internal/models"

	"github.com/rs/zerolog/log"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) RegisterUser(username string, signature, pubKey []byte) error {
	log.Debug().Msgf("verifying signature for: %s", username)

	if !crypto.Verify(pubKey, []byte(username), signature) {
		log.Debug().Msgf("signature is not valid: %x, %s, %x", pubKey, username, signature)
		return errors.New("signature is not valid")
	}

	if user, err := s.repo.GetUserFromServer(username); err == nil {
		if bytes.Compare(user.PubKey, pubKey) != 0 {
			log.Debug().Msgf("user already exist: %s", username)
			return errors.New("username has been already taken")
		}

		return nil
	}

	log.Debug().Msgf("registering user: %s", username)

	return s.repo.AddUserToServer(&models.User{
		Username: username,
		PubKey:   pubKey,
	})
}
