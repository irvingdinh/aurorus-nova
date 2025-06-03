package cmd

import (
	"log"

	"github.com/pocketbase/pocketbase"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Lorem ipsum dolor sit amet",
	Run: func(cmd *cobra.Command, args []string) {
		app := pocketbase.New()

		if err := app.Start(); err != nil {
			log.Fatal(err)
		}
	},
}
