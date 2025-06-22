package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	user "github.com/alanchchen/go-project-skeleton/pkg/user/api"
)

func init() {
	RootCmd.AddCommand(findUserCommand)

	findUserCommand.Flags().String("name", "", "the user name")
	findUserCommand.Flags().Int64("id", -1, "the user id")
}

var findUserCommand = &cobra.Command{
	Use:   "find",
	Short: "finds an new user by ID or name",
	Long:  "finds an new user by ID or name",
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

		if name := viper.GetString("name"); name != "" {
			resp, err = client.FindUserByName(ctx, &user.FindUserByNameRequest{
				Name: name,
			})
		} else if id := viper.GetInt64("id"); id >= 0 {
			resp, err = client.FindUserById(ctx, &user.FindUserByIdRequest{
				Id: id,
			})
		}
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
