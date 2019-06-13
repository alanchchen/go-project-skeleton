package main

import (
	"context"
	"fmt"

	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/grpc"

	"github.com/alanchchen/go-project-skeleton/pkg/api/greeter"
	"github.com/alanchchen/go-project-skeleton/pkg/app"
)

func init() {
	rootCmd.AddCommand(clientCmd)

	clientCmd.Flags().String("name", "", "Tell the server who you are")
}

var clientCmd = &app.Command{
	Use:   "client",
	Short: "client is a greeter client",
	Long:  "client is a greeter client",
	RunE: func(cmd *app.Command, args []string) error {
		return app.RunCustom(cmd, args, func(cfg EndpointConfig, appCfg *app.Config) error {
			conn, err := grpc.Dial(cfg.Endpoint(), grpc.WithInsecure())
			if err != nil {
				return err
			}

			client := greeter.NewServiceClient(conn)

			resp, err := client.SayHello(context.Background(), &greeter.HelloRequest{
				Name: appCfg.GetString("name"),
			})
			if err != nil {
				return err
			}

			fmt.Println("Server says:", resp.Message)
			return nil
		})
	},
}
