package cmd

import (
	"github.com/manthan307/corebase/db"
	"github.com/manthan307/corebase/server"
	"github.com/manthan307/corebase/utils/configs"
	"github.com/manthan307/corebase/utils/logger"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

var (
	Port        int = 8000
	EnvFilePath string
	serveCmd    = &cobra.Command{
		Use:   "serve",
		Short: "serve the application i guess",
		Run: func(cmd *cobra.Command, args []string) {
			if EnvFilePath != "" {
				configs.LoadEnv(EnvFilePath)
			}

			fx.New(
				fx.WithLogger(func() fxevent.Logger {
					return &fxevent.ZapLogger{Logger: logger.ProvideFxLogger()}
				}),
				fx.Provide(logger.ProvideAppLogger),
				db.Module,
				fx.Provide(LoadConfig),
				server.Module,
			).Run()
		},
	}
)

func LoadConfig(client *db.Client) configs.Config {
	return configs.InitConfig(Port, client)
}

func init() {
	serveCmd.PersistentFlags().IntVarP(&Port, "port", "p", 8000, "To set the application port")
	serveCmd.PersistentFlags().StringVarP(&EnvFilePath, "env", "e", "", "To set the env file path")
	AddCommand(serveCmd)
}
