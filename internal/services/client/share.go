package client

import (
	"context"
	"github.com/google/uuid"
	"github.com/mbrostami/goshare/api/grpc"
	"github.com/mbrostami/goshare/api/grpc/pb"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"io"
	"os"
	"path/filepath"
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

	connections := make([]*grpc.Client, len(servers))
	for i, server := range servers {
		c, err := grpc.NewClient(ctx, server)
		if err != nil {
			return err
		}
		connections[i] = c
	}
	log.Debug().Msg("connection to servers was successful!")

	err = connections[0].ShareInit(ctx, uid, filepath.Base(filePath))
	log.Debug().Msgf("start initializing share: got %+v", err)
	if err != nil {
		log.Error().Msgf("couldn't initialize share %+v", err)
		return err
	}

	chunkChannel := make(chan *pb.ShareRequest)
	eg, _ := errgroup.WithContext(ctx)
	for i, _ := range servers {
		index := i
		eg.Go(func() error {
			if err = connections[index].Share(ctx, uid, chunkChannel); err != nil {
				log.Error().Err(err).Send()
				return err
			}
			log.Debug().Msgf("sharing with server %d finished!", index)
			return nil
		})
	}

	buf := make([]byte, 1024)
	var seq int64
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
		chunkChannel <- &r
	}
	close(chunkChannel)
	return eg.Wait()
}
