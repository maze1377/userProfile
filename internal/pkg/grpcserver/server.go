package grpcserver

import (
	"fmt"
	"net"
	"userProfile/pkg/errors"
	"userProfile/pkg/userProfile"

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpclogrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpcprometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Server struct {
	listener net.Listener
	server   *grpc.Server
}

func New(UserProfileViewServer userProfile.UserProfileServer, logger *logrus.Logger, listenPort int) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", listenPort))
	if err != nil {
		return nil, errors.Wrap(err, "failed to listen")
	}

	logEntry := logger.WithFields(map[string]interface{}{
		"app": "userProfile",
	})

	interceptors := []grpc.UnaryServerInterceptor{
		grpclogrus.UnaryServerInterceptor(logEntry),
		grpcprometheus.UnaryServerInterceptor,
		grpcrecovery.UnaryServerInterceptor(),
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(interceptors...)))
	userProfile.RegisterUserProfileServer(grpcServer, UserProfileViewServer)

	log.Infof("userProfile Server Listening on https://0.0.0.0:%d", listenPort)

	return &Server{
		listener: listener,
		server:   grpcServer,
	}, nil
}

func (s *Server) Serve() error {
	if err := s.server.Serve(s.listener); err != nil {
		return errors.Wrap(err, "failed to serve")
	}
	return nil
}

func (s *Server) Stop() {
	s.server.GracefulStop()
}
