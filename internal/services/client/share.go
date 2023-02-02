package client

import (
	"context"
	"io"
	"os"

	"github.com/google/uuid"
	"github.com/mbrostami/goshare/api/grpc"
	"github.com/mbrostami/goshare/api/grpc/pb"
	"github.com/rs/zerolog/log"
)

func (s *Service) Share(ctx context.Context, filePath string, uid uuid.UUID, servers []string) error {
	file, err := os.Open(filePath)
	if err != nil {
		log.Error().Msgf("failed to open file: %v", err)
		return err
	}
	defer file.Close()

	client, err := grpc.NewClientLoadBalancer(servers)
	if err != nil {
		return err
	}

	log.Debug().Msg("connection to servers was successful!")

	chunkChannel := make(chan *pb.ShareRequest)

	buf := make([]byte, 1024)
	var seq int64
	go func() {
		for {
			seq++
			n, err := file.Read(buf)
			if err == io.EOF {
				log.Debug().Msg("sending chunk to channel finished!")
				close(chunkChannel)
				break
			}
			if err != nil {
				log.Error().Msgf("failed to read file: %v", err)
				close(chunkChannel)
				break
			}
			log.Debug().Msgf("sending chunk to channel seq: %d", seq)
			chunkChannel <- &pb.ShareRequest{
				Identifier:     uid.String(),
				FileName:       filePath,
				SequenceNumber: seq,
				Data:           buf[:n],
			}
		}
	}()

	return client.Share(ctx, uid, chunkChannel)
}
