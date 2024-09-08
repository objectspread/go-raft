package app

import (
	"fmt"

	"github.com/objectspread/go-raft/cmd/server/app/flags"
	"github.com/objectspread/go-raft/cmd/server/app/handler"
	"github.com/objectspread/go-raft/cmd/server/app/server"
	"github.com/objectspread/go-raft/proto-gen/server/api_v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type RaftServer struct {
	logger *zap.Logger
	// hServer    *http.Server
	grpcServer *grpc.Server

	api_v1.UnimplementedPingPongServiceServer
}

// RaftServerParams to construct a new Raft Server.
type RaftServerParams struct {
	Logger *zap.Logger
}

func New(params *RaftServerParams) *RaftServer {
	return &RaftServer{
		logger: params.Logger,
	}
}

func (s *RaftServer) Start(options *flags.RaftServerOptions) error {
	grpcServer, err := server.StartGRPCServer(&server.GRPCServerParams{
		Logger:                  s.logger,
		HostPort:                options.GRPC.HostPort,
		MaxReceiveMessageLength: options.GRPC.MaxReceiveMessageLength,
		MaxConnectionAge:        options.GRPC.MaxConnectionAge,
		MaxConnectionAgeGrace:   options.GRPC.MaxConnectionAgeGrace,
		Handler:                 handler.NewGRPCHandler(s.logger),
	})
	if err != nil {
		return fmt.Errorf("could not start grpc server %w", err)
	}
	s.grpcServer = grpcServer
	return nil
}

func (s *RaftServer) Stop() error {
	s.grpcServer.Stop()
	return nil
}
