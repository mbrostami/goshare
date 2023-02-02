package client

import (
	"context"

	"github.com/google/uuid"
	"github.com/mbrostami/goshare/api/grpc"
	"github.com/rs/zerolog/log"
)

func (s *Service) Receive(ctx context.Context, uid uuid.UUID, servers []string) (string, error) {
	client, err := grpc.NewClientLoadBalancer(servers)
	if err != nil {
		return "", err
	}

	log.Debug().Msg("connection to servers was successful!")

	return client.Receive(ctx, uid)
}
