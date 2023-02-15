package grpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/mbrostami/goshare/pkg/tracer"
	"io"

	"github.com/google/uuid"
	"github.com/mbrostami/goshare/api/grpc/pb"
	"github.com/rs/zerolog/log"
)

const MaxConcurrent = 20

func (c *Client) ShareInit(ctx context.Context, uid uuid.UUID, fileName string, fileSize int64) error {
	ctx, span := tracer.NewSpan(ctx, "sender")
	defer span.End()

	res, err := c.conn.ShareInit(ctx, &pb.ShareInitRequest{
		Identifier: uid.String(),
		FileName:   fileName,
		FileSize:   fileSize,
	})
	if err != nil {
		return err
	}
	if res.Error != "" {
		if res.Message == "retry" {
			return c.ShareInit(ctx, uid, fileName, fileSize)
		}
		return errors.New(res.Error)
	}

	return nil
}

func (c *Client) Share(ctx context.Context, uid uuid.UUID, chunks chan *pb.ShareRequest) error {
	ctx, span := tracer.NewSpan(ctx, "sender")
	defer span.End()

	stream, err := c.conn.Share(ctx)
	if err != nil {
		log.Error().Msgf("failed to stream file: %v", err)
		return err
	}

	for chunk := range chunks {
		log.Debug().Msgf("streaming chunk to server: %d", chunk.SequenceNumber)
		err = stream.Send(chunk)
		if err != nil {
			log.Error().Msgf("failed to send chunk: %v", err)
			return err
		}

		response, err := stream.Recv()

		if err == io.EOF {
			log.Debug().Msg("streaming response finished")
			return nil
		}

		if err != nil {
			return fmt.Errorf("failed to receive response chunk: %v", err)
		}

		log.Debug().Msgf("streaming response received %+v", response)
		if response.Error != "" {
			// TODO retry
			log.Error().Msgf("received response %s", response.Error)
			return fmt.Errorf("received response %s", response.Error)
		}
	}

	return nil
}
