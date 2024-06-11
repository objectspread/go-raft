package app

import (
	"context"
	"fmt"
	"log"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	pb "github.com/objectspread/go-raft/proto-gen/server"
)

type RaftServer struct {
	logger *zap.Logger
	// hServer    *http.Server
	grpcServer *grpc.Server
	pb.UnimplementedHelloWorldServiceServer
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

func (s *RaftServer) SayHello(ctx context.Context, in *pb.HelloWorldRequest) (*pb.HelloWorldResponse, error) {
	return &pb.HelloWorldResponse{Message: "Hello, World! "}, nil
}

func (s *RaftServer) Start() error {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen on port 50051: %v", err)
	}

	s.grpcServer = grpc.NewServer()
	pb.RegisterHelloWorldServiceServer(s.grpcServer, &RaftServer{})
	log.Printf("gRPC server listening at %v", lis.Addr())
	if err := s.grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to listen on port 50051: %v", zap.Error(err))
	}
	return nil
}

func (s *RaftServer) Stop() error {
	s.grpcServer.Stop()
	return nil
}
