package client

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/mbrostami/goshare/api/grpc"
	"github.com/mbrostami/goshare/api/grpc/pb"
	"github.com/mbrostami/goshare/pkg/mempage"
	"github.com/rs/zerolog/log"
	"os"
	"sync"
)

func (s *Service) Receive(ctx context.Context, uid uuid.UUID, servers []string) (string, error) {
	log.Trace().Msg("connection to servers was successful!")

	var fileName string
	var fileSize int64
	resChan := make(chan *pb.ReceiveResponse)

	connections := make([]*grpc.Client, 0, len(servers))
	var err error
	for _, server := range servers {
		connection, err := grpc.NewClient(ctx, server)
		if err != nil {
			return "", err
		}

		connections = append(connections, connection)

		res, fs, err := connection.ReceiveInit(ctx, uid)
		log.Trace().Msgf("receive initialize: %s: %v: %+v", server, res, err)
		if res != "" {
			fileName = res
			fileSize = fs
		}
	}

	if fileName == "" || fileSize < 1 {
		return "", errors.New("couldn't get fileName")
	}

	log.Info().Msgf("receiving file %s size: %d", fileName, fileSize)

	wg := sync.WaitGroup{}
	for i, connection := range connections {
		wg.Add(1)
		go func(conn *grpc.Client, index int) {
			if err = conn.Receive(ctx, uid, resChan); err != nil {
				log.Error().Err(err).Send()
			}
			log.Trace().Msgf("wg done! %d", index)
			wg.Done()
		}(connection, i)
	}

	go func() {
		wg.Wait()
		resChan <- &pb.ReceiveResponse{
			SequenceNumber: -1,
		}
	}()

	// blocked by chanel
	if err = s.writeToFile(fileName, resChan); err != nil {
		return "", err
	}

	return fileName, nil
}

func (s *Service) writeToFile(fileName string, resChan chan *pb.ReceiveResponse) error {
	mem := mempage.New()

	// Create a file to store the received chunks
	file, err := os.Create(fileName)
	if err != nil {
		log.Error().Err(err).Send()
		return err
	}
	defer file.Close()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		readChan := make(chan *mempage.Element)
		mem.Read(readChan)
		for elem := range readChan {
			if _, err := file.Write(elem.Data); err != nil {
				log.Error().Err(err).Send()
			}
		}
		wg.Done()
	}()
	for res := range resChan {
		if res.SequenceNumber < 0 {
			log.Trace().Msgf("received %+v", res)
			break
		}
		mem.Store(&mempage.Element{
			Sequence: res.SequenceNumber,
			Data:     res.Data,
		})
	}
	mem.Close()
	wg.Wait()
	return nil
}
