package client

import (
	"context"
	"github.com/google/uuid"
	"github.com/mbrostami/goshare/api/grpc"
	"github.com/mbrostami/goshare/api/grpc/pb"
	"github.com/mbrostami/goshare/pkg/tracer"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"io"
	"os"
	"path/filepath"
)

const (
	KB = 1 << (10 * (iota + 1))
	MB
	GB
)
const MaxConcurrentShare = 0
const ChunkSize = 1 * MB

func (s *Service) Share(ctx context.Context, filePath string, uid uuid.UUID, servers []string) error {
	ctx, span := tracer.NewSpan(ctx, "sender-service")
	defer span.End()

	file, err := os.Open(filePath)
	if err != nil {
		log.Error().Msgf("failed to open file: %v", err)
		return err
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		return err
	}

	connections := make([]*grpc.Client, len(servers))
	for i, server := range servers {
		c, err := grpc.NewClient(ctx, server)
		if err != nil {
			return err
		}
		connections[i] = c
	}
	log.Debug().Msg("connection to servers was successful!")

	for i, _ := range servers {
		err = connections[i].ShareInit(ctx, uid, filepath.Base(filePath), fi.Size())
		log.Debug().Msgf("start initializing share with server %d: got %+v", i, err)
		if err != nil {
			log.Error().Msgf("couldn't initialize share %+v", err)
			return err
		}
	}

	// chunkChannel := make(chan *pb.ShareRequest, fi.Size()%1024)
	chunkChannel := make(chan *pb.ShareRequest, len(servers)*2) // 2 per server

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

	var breakLoop bool
	go func() {
		if err := eg.Wait(); err != nil {
			breakLoop = true
		}
	}()

	buf := make([]byte, ChunkSize) // TODO negotiate with receiver to set the chunk size
	var seq int64
	for {
		if breakLoop {
			break
		}
		seq++

		n, err := file.Read(buf)
		if err == io.EOF {
			log.Debug().Msg("sending chunk to channel finished!")
			// send nil to receiver so receiver knows it's done
			chunkChannel <- &pb.ShareRequest{
				Identifier:     uid.String(),
				SequenceNumber: -1,
			}
			break
		}
		if err != nil {
			log.Error().Msgf("failed to read file: %v", err)

			// send nil to receiver so receiver knows it's done
			chunkChannel <- &pb.ShareRequest{
				Identifier:     uid.String(),
				SequenceNumber: -1,
			}
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
