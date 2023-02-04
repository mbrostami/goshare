package client

import (
	"context"
	"github.com/google/uuid"
	"github.com/mbrostami/goshare/api/grpc"
	"github.com/mbrostami/goshare/api/grpc/pb"
	"github.com/rs/zerolog/log"
	"os"
	"sync"
)

func (s *Service) Receive(ctx context.Context, uid uuid.UUID, servers []string) (string, error) {
	log.Debug().Msg("connection to servers was successful!")

	resChan := make(chan *pb.ReceiveResponse)
	wg := &sync.WaitGroup{}
	for _, server := range servers {
		c, err := grpc.NewClient(ctx, server)
		if err != nil {
			return "", err
		}
		wg.Add(1)
		go func() {
			if err = c.Receive(ctx, uid, resChan); err != nil {
				log.Error().Err(err).Send()
			}
			wg.Done()
		}()
	}

	var fileName string
	var file *os.File
	defer file.Close()
	var err error

	go func() {
		wg.Wait()
		close(resChan)
	}()

	for res := range resChan {
		if res.SequenceNumber < 0 {
			log.Debug().Msgf("wg done!")
			close(resChan)
		}
		if fileName == "" {
			fileName = res.FileName
			// Create a file to store the received chunks
			log.Debug().Msgf("creating the file %s", fileName)
			file, err = os.Create(fileName)
			if err != nil {
				log.Error().Err(err).Send()
				return "", err
			}
		}

		log.Debug().Msgf("writing seq: %d", res.SequenceNumber)

		if _, err = file.Write(res.Data); err != nil {
			log.Error().Err(err).Send()
			return "", err
		}
	}

	wg.Wait()
	return fileName, nil
}
