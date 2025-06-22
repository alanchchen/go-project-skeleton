package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/alanchchen/go-project-skeleton/pkg/app"
)

func init() {
	// Setup flags
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", fmt.Sprintf("config file (default is ./%s.yaml)", appName))

	RootCmd.PersistentFlags().String("api.host", "", "the grpc server listening host")
	RootCmd.PersistentFlags().Int("api.port", 8088, "the grpc server listening port")
}

const (
	appName = "user"
)

var (
	cfgFile string
)

var RootCmd = &cobra.Command{
	Use:   appName,
	Short: appName + " is an RPC server",
	Long:  appName + " is an RPC server",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		app.BindFlags(cmd)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}
