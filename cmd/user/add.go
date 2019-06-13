package main

import (
	"context"
	"encoding/json"
	"fmt"

	_ "github.com/joho/godotenv/autoload"

	"github.com/alanchchen/go-project-skeleton/pkg/api/user"
	"github.com/alanchchen/go-project-skeleton/pkg/app"
)

func init() {
	rootCmd.AddCommand(addUserCommand)

	addUserCommand.Flags().String("name", "", "the user name")
}

var addUserCommand = &app.Command{
	Use:   "add",
	Short: "adds an new user",
	Long:  "adds an new user",
	RunE: func(cmd *app.Command, args []string) error {
		initializers := []interface{}{
			NewConnection,
			NewClient,
		}

		return app.RunCustom(cmd, args, func(client user.ServiceClient, cfg *app.Config) error {
			resp, err := client.AddUser(context.Background(), &user.AddUserRequest{
				Name: cfg.GetString("name"),
			})
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
