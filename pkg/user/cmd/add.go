package cmd

import (
	"context"
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
	RootCmd.AddCommand(addUserCommand)

	addUserCommand.Flags().String("name", "Adam", "the user name")
}

var addUserCommand = &cobra.Command{
	Use:   "add",
	Short: "adds an new user",
	Long:  "adds an new user",
	RunE: func(cmd *cobra.Command, args []string) error {
		host := viper.GetString("api.host")
		port := viper.GetInt("api.port")

		username := viper.GetString("name")

		conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", host, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return err
		}

		ctx, cancel := context.WithDeadlineCause(context.Background(), time.Now().Add(5*time.Second), errors.New("timeout"))
		defer cancel()

		client := user.NewUserServiceClient(conn)
		_, err = client.AddUser(ctx, &user.AddUserRequest{
			Name: username,
		})
		if err != nil {
			return err
		}

		fmt.Printf("User '%s' added successfully\n", username)
		return nil
	},
}
