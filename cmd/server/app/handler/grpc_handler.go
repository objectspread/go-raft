package handler

import (
	"context"
	"fmt"

	"github.com/objectspread/go-raft/proto-gen/server/api_v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type GRPCHandler struct {
	Logger *zap.Logger
	api_v1.UnimplementedPingPongServiceServer
}

// NewGRPCHandler registers routes for this handler on the given router.
func NewGRPCHandler(logger *zap.Logger) *GRPCHandler {
	return &GRPCHandler{
		Logger: logger,
	}
}

func (gh *GRPCHandler) RegisterHandlers(server *grpc.Server) {
	api_v1.RegisterPingPongServiceServer(server, gh)
}

// Ping implements gRPC PingService.
func (gh *GRPCHandler) Ping(ctx context.Context, r *api_v1.PingRequest) (*api_v1.PongResponse, error) {
	var err error
	fmt.Println(ctx, r)
	return &api_v1.PongResponse{
		Message: "Pong",
	}, err
}
