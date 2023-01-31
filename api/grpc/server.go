package grpc

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/mbrostami/goshare/internal/services/server"

	"github.com/rs/zerolog/log"

	"google.golang.org/grpc/keepalive"

	"google.golang.org/grpc"

	"github.com/mbrostami/goshare/api/grpc/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	serverService *server.Service
	pb.UnimplementedGoShareServer
}

func NewServer(serverService *server.Service) *Server {
	return &Server{
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

func (s *Server) InitShare(ctx context.Context, req *pb.InitShareRequest) (*pb.InitShareResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InitShare not implemented")
}

func (s *Server) Share(server pb.GoShare_ShareServer) error {
	return status.Errorf(codes.Unimplemented, "method Share not implemented")
}

func (s *Server) Receive(req *pb.ReceiveRequest, receiver pb.GoShare_ReceiveServer) error {
	return status.Errorf(codes.Unimplemented, "method Receive not implemented")
}

func (s *Server) Register(ctx context.Context, req *pb.RegistrationRequest) (*pb.RegistrationResponse, error) {
	err := s.serverService.RegisterUser(req.Username, req.Signature, req.PubKey)
	if err != nil {
		return &pb.RegistrationResponse{
			Error: fmt.Sprintf("register user: %v", err),
		}, nil
	}

	return &pb.RegistrationResponse{
		Message: "username registered!",
	}, nil
}
