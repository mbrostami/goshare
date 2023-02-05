package client

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/mbrostami/goshare/api/grpc"
	"github.com/mbrostami/goshare/api/grpc/pb"
	"github.com/rs/zerolog/log"
	"os"
	"sync"
)

func (s *Service) Receive(ctx context.Context, uid uuid.UUID, servers []string) (string, error) {
	log.Debug().Msg("connection to servers was successful!")

	var fileName string
	resChan := make(chan *pb.ReceiveResponse)
	wg := &sync.WaitGroup{}
	initwq := &sync.WaitGroup{}
	initctx, cancel := context.WithCancel(ctx)
	defer cancel()
	for _, server := range servers {
		c, err := grpc.NewClient(ctx, server)
		if err != nil {
			return "", err
		}
		initwq.Add(1)
		go func() {
			res, e := c.ReceiveInit(initctx, uid)
			log.Debug().Msgf("receive initialize: %s: %v: %+v", server, res, e)
			if res != "" {
				fileName = res
				cancel()
			}
			initwq.Done()
		}()

		wg.Add(1)
		go func() {
			if err = c.Receive(ctx, uid, resChan); err != nil {
				log.Error().Err(err).Send()
			}
			wg.Done()
		}()
	}

	initwq.Wait()

	if fileName == "" {
		return "", errors.New("couldn't get fileName")
	}

	var file *os.File
	var err error
	// Create a file to store the received chunks
	log.Debug().Msgf("creating the file %s", fileName)
	file, err = os.Create(fileName)
	defer file.Close()
	if err != nil {
		log.Error().Err(err).Send()
		return "", err
	}

	go func() {
		wg.Wait()
		// release the channel when s.Receive is done
		close(resChan)
	}()

	// blocked by chanel
	for res := range resChan {
		if res.SequenceNumber < 0 {
			log.Debug().Msgf("wg done!")
			continue
		}

		log.Trace().Msgf("writing seq: %d", res.SequenceNumber)

		if _, err = file.Write(res.Data); err != nil {
			log.Error().Err(err).Send()
			return "", err
		}
	}

	return fileName, nil
}
