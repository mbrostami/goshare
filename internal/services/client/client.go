package client

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/mbrostami/goshare/api/grpc"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) VerifyServers(ctx context.Context, servers []string) error {
	for _, server := range servers {
		dialCtx, _ := context.WithTimeout(context.Background(), 15*time.Second)
		c, err := grpc.NewClient(dialCtx, server)
		if err != nil {
			return err
		}
		if err = c.Ping(dialCtx); err != nil {
			return fmt.Errorf("server %s is not responding: %v", server, err)
		}
	}

	return nil
}

func (s *Service) GenerateKey(servers []string) (string, uuid.UUID) {
	id := uuid.New()
	str := fmt.Sprintf("%s/%s", id.String(), strings.Join(servers, "-"))
	return base64.StdEncoding.EncodeToString([]byte(str)), id
}

func (s *Service) ParseKey(key string) (servers []string, id uuid.UUID, err error) {
	str, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return
	}

	parts := strings.Split(string(str), "/")
	if len(parts) != 2 {
		err = errors.New("key is not valid")
		return
	}

	id, err = uuid.Parse(parts[0])
	if err != nil {
		return
	}

	servers = strings.Split(parts[1], "-")
	return
}
