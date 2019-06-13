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
	rootCmd.AddCommand(findUserCommand)

	findUserCommand.Flags().String("name", "", "the user name")
	findUserCommand.Flags().Int64("id", -1, "the user id")
}

var findUserCommand = &app.Command{
	Use:   "find",
	Short: "finds an new user by ID or name",
	Long:  "finds an new user by ID or name",
	RunE: func(cmd *app.Command, args []string) error {
		initializers := []interface{}{
			NewConnection,
			NewClient,
		}

		return app.RunCustom(cmd, args, func(client user.ServiceClient, cfg *app.Config) error {
			var resp *user.Users
			var err error

			if name := cfg.GetString("name"); name != "" {
				resp, err = client.FindUserByName(context.Background(), &user.FindUserByNameRequest{
					Name: name,
				})
			}
			if id := cfg.GetInt64("id"); id >= 0 {
				resp, err = client.FindUserById(context.Background(), &user.FindUserByIdRequest{
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
		}, initializers...)
	},
}
