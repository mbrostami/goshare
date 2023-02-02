package grpc

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/mbrostami/goshare/internal/services/server"

	"github.com/rs/zerolog/log"

	"google.golang.org/grpc/keepalive"

	"google.golang.org/grpc"

	"github.com/mbrostami/goshare/api/grpc/pb"
)

type Server struct {
	mu            sync.RWMutex
	relay         map[string]chan *pb.ReceiveResponse
	serverService *server.Service
	pb.UnimplementedGoShareServer
}

func NewServer(serverService *server.Service) *Server {
	return &Server{
		mu:            sync.RWMutex{},
		relay:         make(map[string]chan *pb.ReceiveResponse),
		serverService: serverService,
	}
}

func ListenAndServe(serverService *server.Service, addr string) error {
	sv := NewServer(serverService)
	s := grpc.NewServer(grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle: 1 * time.Minute,
	}))
	pb.RegisterGoShareServer(s, sv)

	listener, err := net.Listen("tcp", addr)

	if err != nil {
		log.Error().Err(err).Send()
		return err
	}

	log.Info().Msgf("listening on addr: %s", addr)

	return s.Serve(listener)
}

func (s *Server) Share(stream pb.GoShare_ShareServer) error {
	log.Debug().Msg("received new share stream")

	for {
		chunk, err := stream.Recv()

		if err == io.EOF || chunk == nil {
			// process buf as a whole file
			log.Debug().Msg("receiving chunks finished")
			break
		}
		if err != nil {
			log.Error().Err(err).Send()
			stream.Send(&pb.ShareResponse{
				Error: err.Error(),
			})
			break
		}

		log.Debug().Msgf("received chunk %+v", chunk)

		relayCounter := 0
	Relay:
		s.mu.RLock()
		recChan, ok := s.relay[chunk.Identifier]
		s.mu.RUnlock()
		if !ok {
			log.Info().Msg("no receiver found! waiting for 5 seconds...")
			relayCounter++
			if relayCounter > 3 {
				log.Error().Msgf("receiver didn't start receiving")
				err = fmt.Errorf("receiver didn't start receiving %s", chunk.Identifier)
				stream.Send(&pb.ShareResponse{
					Error: err.Error(),
				})
				break
			}
			time.Sleep(5 * time.Second)
			goto Relay
		}

		// TODO send to the receiver channel
		recChan <- &pb.ReceiveResponse{
			FileName:       chunk.FileName,
			SequenceNumber: chunk.SequenceNumber,
			Data:           chunk.Data,
		}
		err = stream.Send(&pb.ShareResponse{
			Message: "ok",
		})
		if err != nil {
			log.Error().Err(err).Send()
			break
		}
	}
	return nil
}

func (s *Server) Receive(req *pb.ReceiveRequest, receiver pb.GoShare_ReceiveServer) error {
	log.Debug().Msgf("receiver start receiving on %s", req.Identifier)

	s.mu.Lock()
	s.relay[req.Identifier] = make(chan *pb.ReceiveResponse)
	s.mu.Unlock()

	select {
	case response := <-s.relay[req.Identifier]:
		log.Debug().Msgf("sending data to receiver %s", req.Identifier)

		err := receiver.Send(response)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Error().Err(err).Send()
		}
		//case <-time.After(30 * time.Second):
		//	break
	}

	close(s.relay[req.Identifier])

	s.mu.Lock()
	delete(s.relay, req.Identifier)
	s.mu.Unlock()

	return nil
}

func (s *Server) Ping(ctx context.Context, _ *pb.PingMsg) (*pb.PongMsg, error) {
	return &pb.PongMsg{Pong: true}, nil
}
