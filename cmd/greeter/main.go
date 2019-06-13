package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/alanchchen/go-project-skeleton/pkg/app"
)

func init() {
	// Setup flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", fmt.Sprintf("config file (default is ./%s.yaml)", appName))

	rootCmd.PersistentFlags().String("api.host", "", "the grpc server listening host")
	rootCmd.PersistentFlags().Int("api.port", 8088, "the grpc server listening port")
}

const (
	appName = "greeter"
)

var cfgFile string
var logger = log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)

// rootCmd is the root command
var rootCmd = &cobra.Command{
	Use:   appName,
	Short: appName + " is an RPC server",
	Long:  appName + " is an RPC server",
	RunE: func(cmd *cobra.Command, args []string) error {
		initializers := []interface{}{
			// actors
			app.NewGRPCServerActor,
			app.NewSignalActor,

			// actors' dependencies
			NewRPCServer,
			NewTCPSocket,
			func() app.Logger {
				return logger
			},
		}

		return app.Run(cmd, args, initializers...)
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		logger.Println(err)
		os.Exit(1)
	}
}
