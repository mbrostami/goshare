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

	var receiverIdentifier string
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

		log.Debug().Msgf("received chunk %d", chunk.SequenceNumber)

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

		receiverIdentifier = chunk.Identifier

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
	// close receiver channel
	s.mu.RLock()
	if recChan, ok := s.relay[receiverIdentifier]; ok {
		close(recChan)
	}
	s.mu.RUnlock()

	return nil
}

func (s *Server) Receive(req *pb.ReceiveRequest, receiver pb.GoShare_ReceiveServer) error {
	log.Debug().Msgf("receiver start receiving on %s", req.Identifier)

	s.mu.Lock()
	s.relay[req.Identifier] = make(chan *pb.ReceiveResponse)
	s.mu.Unlock()

	for response := range s.relay[req.Identifier] {

		log.Debug().Msgf("sending data to receiver %s , seq: %d", req.Identifier, response.SequenceNumber)

		err := receiver.Send(response)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Error().Err(err).Send()
		}
	}

	s.mu.Lock()
	delete(s.relay, req.Identifier)
	s.mu.Unlock()

	return nil
}

func (s *Server) Ping(ctx context.Context, _ *pb.PingMsg) (*pb.PongMsg, error) {
	log.Debug().Msg("Pong!")
	return &pb.PongMsg{Pong: true}, nil
}
