package flags

import (
	"flag"
	"fmt"
	"github.com/objectspread/go-raft/pkg/ports"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"time"
)

const (
	flagSuffixHostPort = "host-port"

	flagSuffixHTTPReadTimeout       = "read-timeout"
	flagSuffixHTTPReadHeaderTimeout = "read-header-timeout"
	flagSuffixHTTPIdleTimeout       = "idle-timeout"

	flagSuffixGRPCMaxReceiveMessageLength = "max-message-size"
	flagSuffixGRPCMaxConnectionAge        = "max-connection-age"
	flagSuffixGRPCMaxConnectionAgeGrace   = "max-connection-age-grace"
)

type GRPCOptions struct {
	HostPort string

	MaxReceiveMessageLength int
	// MaxConnectionAge is a duration for the maximum amount of time a connection may exist.
	// See gRPC's keepalive.ServerParameters#MaxConnectionAge.
	MaxConnectionAge time.Duration
	// MaxConnectionAgeGrace is an additive period after MaxConnectionAge after which the connection will be forcibly closed.
	// See gRPC's keepalive.ServerParameters#MaxConnectionAgeGrace.
	MaxConnectionAgeGrace time.Duration
}

type RaftServerOptions struct {
	GRPC GRPCOptions
}

// AddFlags adds flags for CollectorOptions
func AddFlags(flags *flag.FlagSet) {
}

var grpcServerFlagsCfg = serverFlagsConfig{
	prefix: "raft-server",
}

type serverFlagsConfig struct {
	prefix string
}

func (opts *GRPCOptions) initFromViper(v *viper.Viper, _ *zap.Logger, cfg serverFlagsConfig) error {
	opts.HostPort = ports.FormatHostPort(v.GetString(cfg.prefix + "." + flagSuffixHostPort))
	opts.MaxReceiveMessageLength = v.GetInt(cfg.prefix + "." + flagSuffixGRPCMaxReceiveMessageLength)
	opts.MaxConnectionAge = v.GetDuration(cfg.prefix + "." + flagSuffixGRPCMaxConnectionAge)
	opts.MaxConnectionAgeGrace = v.GetDuration(cfg.prefix + "." + flagSuffixGRPCMaxConnectionAgeGrace)

	return nil
}

// InitFromViper initializes CollectorOptions with properties from viper
func (rOpts *RaftServerOptions) InitFromViper(v *viper.Viper, logger *zap.Logger) (*RaftServerOptions, error) {

	//if err := rOpts.HTTP.initFromViper(v, logger, httpServerFlagsCfg); err != nil {
	//	return cOpts, fmt.Errorf("failed to parse HTTP server options: %w", err)
	//}
	//
	if err := rOpts.GRPC.initFromViper(v, logger, grpcServerFlagsCfg); err != nil {
		return rOpts, fmt.Errorf("failed to parse gRPC server options: %w", err)
	}
	return rOpts, nil
}
