package flags

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	flagSuffixGRPCHost = "grpc-host"
)

type GRPCOptions struct {
	Host string
}

type RaftClientOptions struct {
	GRPC GRPCOptions
}

// AddFlags adds flags for CollectorOptions
func AddFlags(flags *flag.FlagSet) {
}

var grpcServerFlagsCfg = serverFlagsConfig{
	prefix: "raft-client",
}

type serverFlagsConfig struct {
	prefix string
}

func (opts *GRPCOptions) initFromViper(v *viper.Viper, _ *zap.Logger, cfg serverFlagsConfig) error {
	opts.Host = v.GetString(cfg.prefix + "." + flagSuffixGRPCHost)

	return nil
}

// InitFromViper initializes CollectorOptions with properties from viper
func (rOpts *RaftClientOptions) InitFromViper(v *viper.Viper, logger *zap.Logger) (*RaftClientOptions, error) {

	//if err := rOpts.HTTP.initFromViper(v, logger, httpServerFlagsCfg); err != nil {
	//	return cOpts, fmt.Errorf("failed to parse HTTP server options: %w", err)
	//}
	//
	if err := rOpts.GRPC.initFromViper(v, logger, grpcServerFlagsCfg); err != nil {
		return rOpts, fmt.Errorf("failed to parse gRPC server options: %w", err)
	}
	return rOpts, nil
}
