package client

import (
	"context"
	"os"
	"sync"

	"github.com/mbrostami/goshare/api/grpc/pb"

	"github.com/google/uuid"
	"github.com/mbrostami/goshare/api/grpc"
	"github.com/rs/zerolog/log"
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
		go func() {
			wg.Add(1)
			if err = c.Receive(ctx, uid, resChan); err != nil {
				log.Error().Err(err).Send()
			}
		}()
	}

	var fileName string
	var file *os.File
	defer file.Close()
	var err error

	defer close(resChan)

	go func() {
		for res := range resChan {
			if res == nil {
				wg.Done()
				log.Debug().Msgf("wg done!")
				continue
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

			log.Debug().Msgf("writing seq: %d", res.SequenceNumber)

			if _, err = file.Write(res.Data); err != nil {
				log.Error().Err(err).Send()
				break
			}
		}
	}()
	wg.Wait()
	return fileName, err
}
