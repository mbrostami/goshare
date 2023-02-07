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
	var fileSize int64
	resChan := make(chan *pb.ReceiveResponse)

	initwq := &sync.WaitGroup{}
	connections := make([]*grpc.Client, len(servers))
	var err error
	for i, server := range servers {
		connections[i], err = grpc.NewClient(ctx, server)
		if err != nil {
			return "", err
		}
		initwq.Add(1)
		go func(index int, server string) {
			res, fs, err := connections[index].ReceiveInit(ctx, uid)
			log.Trace().Msgf("receive initialize: %s: %v: %+v", server, res, err)
			if res != "" {
				fileName = res
				fileSize = fs
			}
			initwq.Done()
		}(i, server)
	}
	initwq.Wait()

	if fileName == "" || fileSize < 1 {
		return "", errors.New("couldn't get fileName")
	}

	log.Info().Msgf("receiving file %s size: %d", fileName, fileSize)

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
	if err = s.writeToFile(file, resChan); err != nil {
		return "", err
	}

	return fileName, nil
}

func (s *Service) writeToFile(file *os.File, resChan chan *pb.ReceiveResponse) error {
	var lastSeq int64
	//writeChan := make(chan []byte)
	//buffer := make([][]byte, 1024)
	for res := range resChan {
		//if {
		//	buffer[res.SequenceNumber%1024] = res.Data
		//}
		//buffer[lastSeq] = res.Data
		if res.SequenceNumber != lastSeq {
			log.Trace().Msgf("mismatch %d :: %d", res.SequenceNumber, lastSeq)
		}
		if res.SequenceNumber < 0 {
			continue
		}
		//log.Trace().Msgf("writing seq: %d", res.SequenceNumber)
		if _, err := file.Write(res.Data); err != nil {
			log.Error().Err(err).Send()
			return err
		}
		lastSeq++
	}
	return nil
}
