package sharing

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/mbrostami/goshare/pkg/tracer"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mbrostami/goshare/api/grpc"
	"golang.org/x/sync/errgroup"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) VerifyServers(ctx context.Context, servers []string, withTLS, skipVerify bool) error {
	ctx, span := tracer.NewSpan(ctx, "verify-servers")
	defer span.End()

	eg := new(errgroup.Group)
	for _, server := range servers {
		server := server // https://golang.org/doc/faq#closures_and_goroutines
		eg.Go(func() error {
			dialCtx, _ := context.WithTimeout(ctx, 15*time.Second)
			c, err := grpc.NewClient(dialCtx, server, withTLS, skipVerify)
			if err != nil {
				return err
			}

			if err = c.Ping(dialCtx); err != nil {
				return fmt.Errorf("server %s is not responding: %v", server, err)
			}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}

func (s *Service) GenerateKey(servers []string) (string, uuid.UUID) {
	id := uuid.New()
	// id := uuid.MustParse("f4128a7e-2b5a-4954-91c6-c595427c435a")
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
