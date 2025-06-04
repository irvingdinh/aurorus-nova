package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/fx"

	"github.com/irvingdinh/aurorus-nova/internal/pb/pb_server"
	"github.com/irvingdinh/aurorus-nova/internal/service/pb_helper_service"
)

func makeServeCommand() *cobra.Command {
	var httpAddr string

	command := &cobra.Command{
		Use:   "serve",
		Short: "Lorem ipsum dolor sit amet",
		Run: func(cmd *cobra.Command, args []string) {
			fx.New(
				fx.Provide(pb_helper_service.NewPbHelperService),
				fx.Provide(pb_server.NewPbServer),
				fx.Invoke(func(
					_ pb_server.PbServer,
				) {
					//
				}),
			).Run()
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
