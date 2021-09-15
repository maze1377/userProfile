package services

import (
	"context"
	"flag"
	"userProfile/pkg/grpc_health_probe"

	"github.com/sirupsen/logrus"
)

type healthCheck struct{}

var backend = flag.String("b", "localhost:10000", "address of userProfile backend")

func NewHealthCheckServiceServer() *healthCheck {
	hc := &healthCheck{}
	return hc
}

func (h healthCheck) Check(ctx context.Context, request *grpc_health_probe.HealthCheckRequest) (*grpc_health_probe.HealthCheckResponse, error) {
	logrus.Debugf("Checking Probe, %s", request.Service)
	// Check grpc request
	if res := checkGetClientInfo(ctx); !res {
		return &grpc_health_probe.HealthCheckResponse{
			Status: grpc_health_probe.HealthCheckResponse_NOT_SERVING,
		}, nil
	}

	return &grpc_health_probe.HealthCheckResponse{
		Status: grpc_health_probe.HealthCheckResponse_SERVING,
	}, nil
}

func checkGetClientInfo(ctx context.Context) bool {
	// todo
	return true
}
