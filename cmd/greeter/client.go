package main

import (
	"context"
	"fmt"

	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	"github.com/alanchchen/go-project-skeleton/pkg/api/greeter"
)

func init() {
	rootCmd.AddCommand(clientCmd)

	clientCmd.Flags().String("name", "", "Tell the server who you are")
}

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "client is a greeter client",
	Long:  "client is a greeter client",
	RunE: func(cmd *cobra.Command, args []string) error {
		conn, err := grpc.Dial(APIEndpoint(), grpc.WithInsecure())
		if err != nil {
			return err
		}

		client := greeter.NewServiceClient(conn)

		resp, err := client.SayHello(context.Background(), &greeter.HelloRequest{
			Name: viper.GetString("name"),
		})
		if err != nil {
			return err
		}

		fmt.Println("Server says:", resp.Message)
		return nil
	},
}
