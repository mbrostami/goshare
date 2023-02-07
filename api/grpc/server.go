package grpc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/mbrostami/goshare/api/grpc/pb"
	"github.com/mbrostami/goshare/internal/services/server"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type Server struct {
	mu            sync.RWMutex
	relay         map[string]chan *pb.ReceiveResponse
	relayInit     map[string]chan *pb.ReceiveInitResponse
	serverService *server.Service
	pb.UnimplementedGoShareServer
}

func NewServer(serverService *server.Service) *Server {
	return &Server{
		mu:            sync.RWMutex{},
		relay:         make(map[string]chan *pb.ReceiveResponse),
		relayInit:     make(map[string]chan *pb.ReceiveInitResponse),
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

func (s *Server) ShareInit(ctx context.Context, req *pb.ShareInitRequest) (*pb.ShareInitResponse, error) {
	relayCounter := 0
Relay:
	s.mu.RLock()
	recInitChan, ok := s.relayInit[req.Identifier]
	s.mu.RUnlock()
	if !ok {
		log.Info().Msg("no receiver initialized! waiting for 5 seconds...")
		relayCounter++
		if relayCounter > 10 {
			log.Error().Msgf("receiver didn't initialize receiving")
			err := fmt.Errorf("receiver didn't initialize receiving %s", req.Identifier)
			return &pb.ShareInitResponse{
				Message: "failed",
				Error:   err.Error(),
			}, nil
		}
		time.Sleep(5 * time.Second)
		goto Relay
	}

	log.Debug().Msg("sending initialize response to receiving channel")

	recInitChan <- &pb.ReceiveInitResponse{FileName: req.FileName, FileSize: req.FileSize}
	return &pb.ShareInitResponse{
		Message: "ok",
	}, nil
}

func (s *Server) ReceiveInit(ctx context.Context, req *pb.ReceiveRequest) (*pb.ReceiveInitResponse, error) {
	log.Debug().Msgf("receiver initialized receiving on %s", req.Identifier)

	s.mu.Lock()
	s.relayInit[req.Identifier] = make(chan *pb.ReceiveInitResponse, 1)
	s.mu.Unlock()

	var res *pb.ReceiveInitResponse
	for response := range s.relayInit[req.Identifier] {
		res = response
		break
	}

	s.mu.Lock()
	close(s.relayInit[req.Identifier])
	s.relay[req.Identifier] = make(chan *pb.ReceiveResponse, req.Semaphore)
	s.mu.Unlock()

	return res, nil
}

func (s *Server) Share(stream pb.GoShare_ShareServer) error {
	log.Debug().Msg("received new share stream")

	var receiverIdentifier string
	for {
		chunk, err := stream.Recv()

		if err == io.EOF || chunk == nil {
			break
		}
		if err != nil {
			log.Error().Err(err).Send()
			stream.Send(&pb.ShareResponse{
				Message: "retry",
				Error:   err.Error(),
			})
			break
		}

		log.Debug().Msgf("received chunk %d", chunk.SequenceNumber)

		s.mu.RLock()
		recChan, ok := s.relay[chunk.Identifier]
		s.mu.RUnlock()
		if !ok {
			log.Info().Msg("no receiver found!...")
			stream.Send(&pb.ShareResponse{
				Message: "retry",
				Error:   err.Error(),
			})
			break
		}

		receiverIdentifier = chunk.Identifier

		// TODO send to the receiver channel
		recChan <- &pb.ReceiveResponse{
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
	s.mu.Lock()
	if recChan, ok := s.relay[receiverIdentifier]; ok {
		close(recChan)
		delete(s.relay, receiverIdentifier)
	}
	s.mu.Unlock()

	return nil
}

func (s *Server) Receive(req *pb.ReceiveRequest, receiver pb.GoShare_ReceiveServer) error {
	log.Debug().Msgf("receiver start receiving on %s", req.Identifier)

	s.mu.RLock()
	recChan, ok := s.relay[req.Identifier]
	s.mu.RUnlock()
	if !ok {
		return errors.New("receiver is not initialized")
	}

	for response := range recChan {

		log.Debug().Msgf("sending data to receiver %s , seq: %d", req.Identifier, response.SequenceNumber)

		err := receiver.Send(response)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Error().Err(err).Send()
			break
		}
	}

	return nil
}

func (s *Server) Ping(ctx context.Context, _ *pb.PingMsg) (*pb.PongMsg, error) {
	log.Debug().Msg("Pong!")
	return &pb.PongMsg{Pong: true}, nil
}
