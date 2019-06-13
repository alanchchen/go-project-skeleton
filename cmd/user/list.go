package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	_ "github.com/joho/godotenv/autoload"

	"github.com/alanchchen/go-project-skeleton/pkg/api/user"
	"github.com/alanchchen/go-project-skeleton/pkg/app"
)

func init() {
	rootCmd.AddCommand(listUsersCommand)
}

var listUsersCommand = &app.Command{
	Use:   "list",
	Short: "list all users",
	Long:  "list all users",
	RunE: func(cmd *app.Command, args []string) error {
		initializers := []interface{}{
			NewConnection,
			NewClient,
		}

		return app.RunCustom(cmd, args, func(client user.ServiceClient) error {
			resp, err := client.ListUsers(context.Background(), &empty.Empty{})
			if err != nil {
				return err
			}

			rawData, err := json.MarshalIndent(resp, "", "  ")
			if err != nil {
				return err
			}

			fmt.Println(string(rawData))

			return nil
		}, initializers...)
	},
}
