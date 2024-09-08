package main

import (
	"fmt"
	"github.com/objectspread/go-raft/cmd/client/app"
	"github.com/objectspread/go-raft/cmd/client/app/flags"
	"go.uber.org/zap"
	"os"

	"github.com/objectspread/go-raft/pkg/config"
	cmdFlags "github.com/objectspread/go-raft/pkg/flags"
	"github.com/objectspread/go-raft/pkg/version"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
	_ "go.uber.org/automaxprocs"
)

func main() {
	svc := cmdFlags.NewService()

	v := viper.New()
	command := &cobra.Command{
		Use:   "go-raft-client",
		Short: "Go Raft Client - sends requests to raft-servers",
		Long:  "Go Raft Client - sends requests to raft-servers",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := svc.Start(v); err != nil {
				return err
			}
			logger := svc.Logger
			clientOpts, err := new(flags.RaftClientOptions).InitFromViper(v, logger)
			if err != nil {
				logger.Fatal("Failed to initialize go-raft-client", zap.Error(err))
			}
			app := app.New(&app.RaftClientParams{Logger: logger})

			if err := app.Start(clientOpts); err != nil {
				fmt.Println(err)
			}
			svc.RunAndThen(func() {
				logger.Info("Closing down the client service")
			})
			return nil
		},
	}

	command.AddCommand(version.Command())

	config.AddFlags(
		v,
		command,
		flags.AddFlags,
		svc.AddFlags,
	)

	if err := command.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
