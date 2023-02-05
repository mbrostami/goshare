package grpc

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/google/uuid"
	"github.com/mbrostami/goshare/api/grpc/pb"
	"github.com/rs/zerolog/log"
)

func (c *Client) ShareInit(ctx context.Context, uid uuid.UUID, fileName string) error {
	res, err := c.conn.ShareInit(ctx, &pb.ShareInitRequest{
		Identifier: uid.String(),
		FileName:   fileName,
	})
	if err != nil {
		return err
	}
	if res.Error != "" {
		if res.Message == "retry" {
			return c.ShareInit(ctx, uid, fileName)
		}
		return errors.New(res.Error)
	}

	return nil
}

func (c *Client) Share(ctx context.Context, uid uuid.UUID, chunks chan *pb.ShareRequest) error {
	stream, err := c.conn.Share(ctx)
	if err != nil {
		log.Error().Msgf("failed to stream file: %v", err)
		return err
	}

	semaphore := make(chan struct{}, 5)
	for chunk := range chunks {
		semaphore <- struct{}{}

		log.Debug().Msgf("streaming chunk to server: %d : %s", chunk.SequenceNumber)
		err = stream.Send(chunk)
		if err != nil {
			log.Error().Msgf("failed to send chunk: %v", err)
			return err
		}

		response, err := stream.Recv()
		<-semaphore

		if err == io.EOF {
			log.Debug().Msg("streaming response finished")
			return nil
		}

		if err != nil {
			return fmt.Errorf("failed to receive response chunk: %v", err)
		}

		log.Debug().Msgf("streaming response received %+v", response)
		if response.Error != "" {
			if response.Message == "retry" {
				log.Debug().Msgf("resending chunk no: %d", chunk.SequenceNumber)
				chunks <- chunk
				continue
			}
			log.Error().Msgf("received response %s", response.Error)
			return fmt.Errorf("received response %s", response.Error)
		}
	}

	return nil
}
