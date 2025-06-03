package cmd

import (
	"log"

	"github.com/pocketbase/pocketbase"
	"github.com/spf13/cobra"
)

func makeServeCommand() *cobra.Command {
	var httpAddr string

	command := &cobra.Command{
		Use:   "serve",
		Short: "Lorem ipsum dolor sit amet",
		Run: func(cmd *cobra.Command, args []string) {
			app := pocketbase.New()

			if err := app.Start(); err != nil {
				log.Fatal(err)
			}
		},
	}

	command.PersistentFlags().StringVar(
		&httpAddr,
		"http",
		"",
		"TCP address to listen for the HTTP server",
	)

	return command
}
