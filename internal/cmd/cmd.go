package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Short: "Lorem ipsum dolor sit amet",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
		os.Exit(0)
	},
}

func Execute() {
	cmd.AddCommand(makeServeCommand())

	if err := cmd.Execute(); err != nil {
		log.Fatal(err.Error())
	}
}
