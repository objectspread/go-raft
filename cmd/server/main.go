package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	_ "go.uber.org/automaxprocs"
	"go.uber.org/zap"

	"github.com/objectspread/go-raft/cmd/server/app"
)

func main() {
	command := &cobra.Command{
		Use:   "goraft-server",
		Short: "Go Raft server receives and processes requests from clients",
		Long:  "Go Raft server receives and processes requests from clients",
		RunE: func(cmd *cobra.Command, args []string) error {
			app := app.New(&app.RaftServerParams{Logger: &zap.Logger{}})

			if err := app.Start(); err != nil {
				fmt.Println(err)
			}
			return nil
		},
	}

	if err := command.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
