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
	log.Trace().Msg("connection to servers was successful!")

	var fileName string
	resChan := make(chan *pb.ReceiveResponse)

	initwq := &sync.WaitGroup{}
	initctx, cancel := context.WithCancel(ctx)
	defer cancel()
	connections := make([]*grpc.Client, len(servers))
	var err error
	for i, server := range servers {
		connections[i], err = grpc.NewClient(ctx, server)
		if err != nil {
			return "", err
		}
		initwq.Add(1)
		go func(index int, server string) {
			res, e := connections[index].ReceiveInit(initctx, uid)
			log.Trace().Msgf("receive initialize: %s: %v: %+v", server, res, e)
			if res != "" {
				fileName = res
				cancel()
			}
			initwq.Done()
		}(i, server)
	}
	initwq.Wait()

	if fileName == "" {
		return "", errors.New("couldn't get fileName")
	}

	wg := &sync.WaitGroup{}
	for i, _ := range servers {
		wg.Add(1)
		go func(index int) {
			if err = connections[index].Receive(ctx, uid, resChan); err != nil {
				log.Error().Err(err).Send()
			}
			log.Trace().Msgf("wg done! %d", index)
			wg.Done()
		}(i)
	}

	var file *os.File
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
