package cmd

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	"github.com/alanchchen/go-project-skeleton/pkg/greeter/api"
	"github.com/alanchchen/go-project-skeleton/pkg/greeter/core"
)

func init() {
	RootCmd.AddCommand(ServerCmd)
}

var ServerCmd = &cobra.Command{
	Use:          "server",
	Short:        "server is a greeter server",
	Long:         "server is a greeter server",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		host := viper.GetString("api.host")
		port := viper.GetInt("api.port")

		l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
		if err != nil {
			return err
		}
		defer l.Close()

		server := grpc.NewServer()
		greeter := core.NewService()
		api.RegisterGreeterServiceServer(server, greeter)

		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

		go func() {
			if sig := <-signalChan; sig != nil {
				fmt.Println("greeter server received signal:", sig)
				server.GracefulStop()
			}
		}()

		fmt.Println("greeter server started at", fmt.Sprintf("%s:%d", host, port))

		if err := server.Serve(l); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			return err
		}

		fmt.Println("greeter server stopped")

		return nil
	},
}
