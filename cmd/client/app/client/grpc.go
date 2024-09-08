package client

import (
	"github.com/objectspread/go-raft/cmd/server/app/handler"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClientParams struct {
	Host           string
	HostPortActual string
	Logger         *zap.Logger
	OnError        func(error)
	Handler        *handler.GRPCHandler
	IsSecure       bool
}

func StartGRPCClient(params *GRPCClientParams) (*grpc.ClientConn, error) {
	var grpcOpts []grpc.DialOption

	if !params.IsSecure {
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	grpcConn, err := grpc.NewClient(params.Host, grpcOpts...)

	if err != nil {
		return nil, err
	}

	return grpcConn, nil
}
