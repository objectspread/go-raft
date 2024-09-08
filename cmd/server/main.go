package main

import (
	"fmt"
	"go.uber.org/zap"
	"os"

	"github.com/objectspread/go-raft/cmd/server/app/flags"
	"github.com/objectspread/go-raft/pkg/config"
	cmdFlags "github.com/objectspread/go-raft/pkg/flags"
	"github.com/objectspread/go-raft/pkg/version"
	"github.com/spf13/viper"

	"github.com/objectspread/go-raft/cmd/server/app"
	"github.com/spf13/cobra"
	_ "go.uber.org/automaxprocs"
)

func main() {
	svc := cmdFlags.NewService()

	v := viper.New()
	command := &cobra.Command{
		Use:   "go-raft-server",
		Short: "Go Raft server receives and processes requests from raft-clients",
		Long:  "Go Raft server receives and processes requests from raft-clients",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := svc.Start(v); err != nil {
				return err
			}
			logger := svc.Logger
			serverOpts, err := new(flags.RaftServerOptions).InitFromViper(v, logger)
			if err != nil {
				logger.Fatal("Failed to initialize collector", zap.Error(err))
			}
			app := app.New(&app.RaftServerParams{Logger: logger})

			if err := app.Start(serverOpts); err != nil {
				fmt.Println(err)
			}
			svc.RunAndThen(func() {
				fmt.Println("Closing down this service")
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
