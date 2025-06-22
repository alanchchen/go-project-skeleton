package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	user "github.com/alanchchen/go-project-skeleton/pkg/user/api"
)

func init() {
	RootCmd.AddCommand(listUsersCommand)
}

var listUsersCommand = &cobra.Command{
	Use:   "list",
	Short: "list all users",
	Long:  "list all users",
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			resp *user.Users
			host = viper.GetString("api.host")
			port = viper.GetInt("api.port")
		)

		conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", host, port),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			return err
		}

		client := user.NewUserServiceClient(conn)

		ctx, cancel := context.WithDeadlineCause(context.Background(), time.Now().Add(5*time.Second), errors.New("timeout"))
		defer cancel()

		resp, err = client.ListUsers(ctx, &empty.Empty{})
		if err != nil {
			return err
		}

		rawData, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return err
		}

		fmt.Println(string(rawData))

		return nil
	},
}
