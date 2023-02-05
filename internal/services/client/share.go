package client

import (
	"context"
	"golang.org/x/sync/errgroup"
	"io"
	"os"
	"path/filepath"
	"sync"

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

	//TODO fixed size
	//fi, err := file.Stat()
	//if err != nil {
	//	return err
	//}
	//chunkChannel := make(chan *pb.ShareRequest, (fi.Size()%1024)+1)
	chunkChannel := make(chan *pb.ShareRequest)
	defer close(chunkChannel)

	eg, _ := errgroup.WithContext(ctx)
	for i, server := range servers {
		c, err := grpc.NewClient(ctx, server)
		if err != nil {
			return err
		}

		if i == 0 {
			err := c.ShareInit(ctx, uid, filepath.Base(filePath))
			log.Debug().Msgf("start initializing share with %s: got %+v", server, err)
			if err != nil {
				log.Error().Msgf("couldn't initialize share %+v", err)
				return err
			}
		}

		eg.Go(func() error {
			if err = c.Share(ctx, uid, chunkChannel); err != nil {
				log.Error().Err(err).Send()
				return err
			}
			return nil
		})
	}
	log.Debug().Msg("connection to servers was successful!")

	buf := make([]byte, 1024)
	var seq int64
	wg := &sync.WaitGroup{}
	semaphore := make(chan struct{}, 4)
	for {
		seq++
		n, err := file.Read(buf)
		if err == io.EOF {
			log.Debug().Msg("sending chunk to channel finished!")
			// send nil to receiver so receiver knows it's done
			pbr := pb.ShareRequest{
				Identifier:     uid.String(),
				SequenceNumber: -1,
			}
			chunkChannel <- &pbr
			break
		}
		if err != nil {
			log.Error().Msgf("failed to read file: %v", err)

			// send nil to receiver so receiver knows it's done
			pbr := pb.ShareRequest{
				Identifier:     uid.String(),
				SequenceNumber: -1,
			}
			chunkChannel <- &pbr
			break
		}
		log.Debug().Msgf("sending chunk to channel seq: %d", seq)
		r := pb.ShareRequest{
			Identifier:     uid.String(),
			SequenceNumber: seq,
		}
		r.Data = make([]byte, n)
		copy(r.Data, buf[:n])
		wg.Add(1)
		semaphore <- struct{}{}
		go func() {
			chunkChannel <- &r
			<-semaphore
			wg.Done()
		}()
	}
	wg.Wait()
	return nil
}
