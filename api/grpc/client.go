package grpc

import (
	"context"
	"errors"

	"github.com/mbrostami/goshare/api/grpc/pb"
	"google.golang.org/grpc"
)

type Client struct {
	conn pb.GoShareClient
}

func NewClient(addr string) (*Client, error) {
	conn, err := grpc.Dial(
		addr,
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	var client Client
	client.conn = pb.NewGoShareClient(conn)

	return &client, nil
}

func (c *Client) Register(ctx context.Context, username string, pubKey, signature []byte) error {
	res, err := c.conn.Register(ctx, &pb.RegistrationRequest{
		Username:  username,
		PubKey:    pubKey,
		Signature: signature,
	})
	if err != nil {
		return err
	}

	if res.Error != "" {
		return errors.New(res.Error)
	}

	return nil
}
