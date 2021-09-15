package main

import (
	"net"

	services "userProfile/internal/pkg/grpc_health_probe"
	"userProfile/pkg/grpc_health_probe"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Liveness struct {
	conn net.Listener
}

func (l *Liveness) startServer() {
	addr := "0.0.0.0:10001"
	conn, err := net.Listen("tcp", addr)
	if err != nil {
		log.WithError(err).Fatalf("Failed to listen on %q", addr)
	}
	l.conn = conn

	healthCheckServiceServer := services.NewHealthCheckServiceServer()
	var opts []grpc.ServerOption
	gRPCServer := grpc.NewServer(opts...)
	grpc_health_probe.RegisterHealthServer(gRPCServer, healthCheckServiceServer)

	log.Infof("Liveness Server Listening on https://%s", addr)
	if err = gRPCServer.Serve(conn); err != nil {
		log.WithError(err).Fatal("Error while apiHandler.Serve()")
	}
}
func (l *Liveness) stopServer() {
	err := l.conn.Close()
	if err != nil {
		log.Error(err)
	}
}
