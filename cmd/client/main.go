package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/objectspread/go-raft/cmd/client/app"
)

func main() {

	command := &cobra.Command{
		Use:   "goraft-client",
		Short: "Go Raft client sends requests to Raft server",
		Long:  `Go Raft client sends requests to Raft server`,
		RunE: func(cmd *cobra.Command, args []string) error {
			app := app.New(&app.RaftClientParams{Logger: &zap.Logger{}})

			if err := app.Connect(); err != nil {
				fmt.Println(err)
			}
			defer app.Close()

			if err := app.SendHelloRequest(); err != nil {
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
