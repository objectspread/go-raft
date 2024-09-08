// Copyright (c) 2019 The Jaeger Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package flags

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapgrpc"
	"google.golang.org/grpc/grpclog"
)

// Service represents an abstract Jaeger backend component with some basic shared functionality.
type Service struct {
	// Logger is initialized after parsing Viper flags like --log-level.
	Logger *zap.Logger

	signalsChannel chan os.Signal
}

// NewService creates a new Service.
func NewService() *Service {
	signalsChannel := make(chan os.Signal, 1)
	signal.Notify(signalsChannel, os.Interrupt, syscall.SIGTERM)

	return &Service{
		signalsChannel: signalsChannel,
	}
}

// AddFlags registers CLI flags.
func (s *Service) AddFlags(flagSet *flag.FlagSet) {
	AddConfigFileFlag(flagSet)
	AddFlags(flagSet)
}

// Start bootstraps the service and starts the admin server.
func (s *Service) Start(v *viper.Viper) error {
	if err := TryLoadConfigFile(v); err != nil {
		return fmt.Errorf("cannot load config file: %w", err)
	}

	sFlags := new(SharedFlags).InitFromViper(v)
	newProdConfig := zap.NewProductionConfig()
	// newProdConfig.Sampling = nil
	logger, err := sFlags.NewLogger(newProdConfig)
	if err != nil {
		return fmt.Errorf("cannot create logger: %w", err)
	}
	s.Logger = logger
	grpclog.SetLoggerV2(zapgrpc.NewLogger(
		logger.WithOptions(
			zap.AddCallerSkip(5), // ensure the actual caller:lineNo is shown
		)))

	return nil
}

// RunAndThen sets the health check to Ready and blocks until SIGTERM is received.
// If then runs the shutdown function and exits.
func (s *Service) RunAndThen(shutdown func()) {
	// s.HC().Ready()

	<-s.signalsChannel

	s.Logger.Info("Shutting down")
	// s.HC().Set(healthcheck.Unavailable)

	if shutdown != nil {
		shutdown()
	}

	// s.Admin.Close()
	s.Logger.Info("Shutdown complete")
}
