package app

import (
	"fmt"
	"github.com/objectspread/go-raft/cmd/client/app/client"
	"github.com/objectspread/go-raft/cmd/client/app/flags"
	"github.com/objectspread/go-raft/cmd/server/app/handler"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type RaftClient struct {
	logger *zap.Logger
	// hServer    *http.Server
	grpcConn *grpc.ClientConn
}

// RaftClientParams to construct a new Raft Client.
type RaftClientParams struct {
	Logger *zap.Logger
	Host   string
}

func New(params *RaftClientParams) *RaftClient {
	return &RaftClient{
		logger: params.Logger,
	}
}

func (c *RaftClient) Start(options *flags.RaftClientOptions) error {
	grpcConn, err := client.StartGRPCClient(&client.GRPCClientParams{
		Logger:  c.logger,
		Host:    options.GRPC.Host,
		Handler: handler.NewGRPCHandler(c.logger),
	})
	if err != nil {
		return fmt.Errorf("could not start grpc server %w", err)
	}

	//_c := api_v1.NewPingPongServiceClient(grpcConn)
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	//defer cancel()
	//
	//r, err := _c.Ping(ctx, &api_v1.PingRequest{})
	//if err != nil {
	//	return err
	//}
	//
	//log.Printf("Response from gRPC server's SayHello function: %s", r.GetMessage())

	c.grpcConn = grpcConn
	return nil
}

func (c *RaftClient) Stop() error {
	c.grpcConn.Close()
	return nil
}
