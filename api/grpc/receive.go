package grpc

import (
	"context"
	"github.com/google/uuid"
	"github.com/mbrostami/goshare/api/grpc/pb"
	"github.com/rs/zerolog/log"
	"io"
	"time"
)

func (c *Client) ReceiveInit(ctx context.Context, id uuid.UUID) (string, error) {
	ctx, _ = context.WithTimeout(ctx, 10*time.Second)
	res, err := c.conn.ReceiveInit(ctx, &pb.ReceiveRequest{
		Identifier: id.String(),
	})
	if err != nil {
		return "", err
	}
	return res.FileName, nil
}

func (c *Client) Receive(ctx context.Context, id uuid.UUID, resChan chan *pb.ReceiveResponse) error {
	log.Debug().Msgf("client sending receive request")
	stream, err := c.conn.Receive(ctx, &pb.ReceiveRequest{Identifier: id.String()})
	if err != nil {
		return err
	}
	defer stream.CloseSend()

	for {
		res, err := stream.Recv()

		log.Debug().Msgf("received from server %+v, err: %+v", res, err)

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Error().Err(err).Send()
			break
		}
		resChan <- res
	}
	return nil
}
