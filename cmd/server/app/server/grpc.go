package server

import (
	"fmt"
	"net"
	"time"

	"github.com/objectspread/go-raft/cmd/server/app/handler"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type GRPCServerParams struct {
	MaxReceiveMessageLength int
	MaxConnectionAge        time.Duration
	MaxConnectionAgeGrace   time.Duration
	HostPort                string
	HostPortActual          string
	Logger                  *zap.Logger
	OnError                 func(error)
	Handler                 *handler.GRPCHandler
}

func StartGRPCServer(params *GRPCServerParams) (*grpc.Server, error) {
	var server *grpc.Server
	var grpcOpts []grpc.ServerOption

	if params.MaxReceiveMessageLength > 0 {
		grpcOpts = append(grpcOpts, grpc.MaxRecvMsgSize(params.MaxReceiveMessageLength))
	}

	server = grpc.NewServer(grpcOpts...)
	reflection.Register(server)

	listener, err := net.Listen("tcp", params.HostPort)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on gRPC port: %w", err)
	}
	params.HostPortActual = listener.Addr().String()

	if err := serveGRPC(server, listener, params); err != nil {
		return nil, err
	}

	return server, nil
}

func serveGRPC(server *grpc.Server, listener net.Listener, params *GRPCServerParams) error {
	healthServer := health.NewServer()

	params.Handler.RegisterHandlers(server)

	healthServer.SetServingStatus("go-raft.api_v2.Server", grpc_health_v1.HealthCheckResponse_SERVING)

	grpc_health_v1.RegisterHealthServer(server, healthServer)
	params.Logger.Info(params.HostPortActual)

	params.Logger.Info("Starting go-raft gRPC server", zap.String("raft-server.host-port", params.HostPortActual))
	go func() {
		if err := server.Serve(listener); err != nil {
			params.Logger.Error("Could not launch gRPC service", zap.Error(err))
			if params.OnError != nil {
				params.OnError(err)
			}
		}
	}()

	return nil
}
