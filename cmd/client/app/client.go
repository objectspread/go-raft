package app

import (
	"context"
	"log"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/objectspread/go-raft/proto-gen/server"
)

type RaftClient struct {
	logger *zap.Logger
	// hServer    *http.Server
	grpcConn *grpc.ClientConn
}

// RaftServerParams to construct a new Raft Server.
type RaftClientParams struct {
	Logger *zap.Logger
	Host   string
}

func New(params *RaftClientParams) *RaftClient {
	d := time.Hour
	return &RaftClient{
		logger: params.Logger,
	}
}

func (c *RaftClient) Close() error {
	return c.grpcConn.Close()
}

func (c *RaftClient) Connect() (err error) {
	c.grpcConn, err = grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	return nil
}

func (c *RaftClient) SendHelloRequest() (err error) {

	_c := pb.NewHelloWorldServiceClient(c.grpcConn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := _c.SayHello(ctx, &pb.HelloWorldRequest{})
	if err != nil {
		return err
	}

	log.Printf("Response from gRPC server's SayHello function: %s", r.GetMessage())
	return nil
}
