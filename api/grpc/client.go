package grpc

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/google/uuid"

	"github.com/mbrostami/goshare/api/grpc/pb"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type Client struct {
	conn pb.GoShareClient
}

func NewClient(ctx context.Context, addr string) (*Client, error) {
	conn, err := grpc.DialContext(
		ctx,
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

func (c *Client) Ping(ctx context.Context) error {
	_, err := c.conn.Ping(ctx, &pb.PingMsg{Ping: true})
	return err
}

func (c *Client) Share(ctx context.Context, uid uuid.UUID, chunks chan *pb.ShareRequest) error {
	stream, err := c.conn.Share(ctx)
	defer stream.CloseSend()
	if err != nil {
		log.Error().Msgf("failed to stream file: %v", err)
		return err
	}

	for chunk := range chunks {
		log.Debug().Msgf("streaming chunk to server: %+v", chunk)
		err := stream.Send(chunk)
		if err != nil {
			log.Error().Msgf("failed to send chunk: %v", err)
			return err
		}
		response, err := stream.Recv()
		if err == io.EOF {
			log.Debug().Msg("streaming response finished")
			break
		}
		if err != nil {
			return fmt.Errorf("failed to receive response chunk: %v", err)
		}

		log.Debug().Msgf("streaming response received %+v", response)
		if response.Error != "" {
			log.Error().Msgf("received response %s", response.Error)
			return fmt.Errorf("received response %s", response.Error)
		}
	}
	return nil
}

func (c *Client) Receive(ctx context.Context, id uuid.UUID, opts ...grpc.CallOption) (string, error) {
	log.Debug().Msgf("client sending receive request")
	stram, err := c.conn.Receive(ctx, &pb.ReceiveRequest{Identifier: id.String()})
	if err != nil {
		return "", err
	}
	defer stram.CloseSend()

	//buf := make([][]byte, 0)
	var fileName string
	var file *os.File
	defer file.Close()
	var nextSequence int64
	nextSequence = 1
	for {
		res, err := stram.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Error().Err(err).Send()
			break
		}
		if fileName == "" {
			fileName = res.FileName
			// Create a file to store the received chunks
			file, err = os.Create(fileName)
			if err != nil {
				log.Error().Err(err).Send()
				break
			}
		}

		if nextSequence != res.SequenceNumber {
			log.Debug().Msgf("sequence not matched %d, %d skipping...", nextSequence, res.SequenceNumber)
			continue
		}

		nextSequence++
		if _, err = file.Write(res.Data); err != nil {
			log.Error().Err(err).Send()
			break
		}
	}
	return fileName, nil
}
