package grpc

import (
	"context"
	"github.com/google/uuid"
	"github.com/mbrostami/goshare/api/grpc/pb"
	"github.com/rs/zerolog/log"
	"io"
)

func (c *Client) ReceiveInit(ctx context.Context, id uuid.UUID) (string, int64, error) {
	res, err := c.conn.ReceiveInit(ctx, &pb.ReceiveRequest{
		Identifier: id.String(),
	})
	if err != nil {
		return "", 0, err
	}
	return res.FileName, res.FileSize, nil
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

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Error().Err(err).Send()
			break
		}

		log.Debug().Msgf("received %d from server", res.SequenceNumber)
		resChan <- res
	}
	return nil
}
