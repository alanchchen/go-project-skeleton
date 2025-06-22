package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	greeter "github.com/alanchchen/go-project-skeleton/pkg/greeter/api"
)

func init() {
	RootCmd.AddCommand(clientCmd)

	clientCmd.Flags().String("name", "", "Tell the server who you are")
}

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "client is a greeter client",
	Long:  "client is a greeter client",
	RunE: func(cmd *cobra.Command, args []string) error {
		host := viper.GetString("api.host")
		port := viper.GetInt("api.port")

		conn, err := grpc.Dial(fmt.Sprintf("%s:%d", host, port), grpc.WithInsecure())
		if err != nil {
			return err
		}

		client := greeter.NewGreeterServiceClient(conn)

		resp, err := client.SayHello(context.Background(), &greeter.HelloRequest{
			Name: cmd.Flag("name").Value.String(),
		})
		if err != nil {
			return err
		}

		fmt.Println("Server says:", resp.Message)
		return nil
	},
}
