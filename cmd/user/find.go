package main

import (
	"context"
	"encoding/json"
	"fmt"

	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/dig"

	"github.com/alanchchen/go-project-skeleton/pkg/api/user"
)

func init() {
	rootCmd.AddCommand(findUserCommand)

	findUserCommand.Flags().String("name", "", "the user name")
	findUserCommand.Flags().Int64("id", -1, "the user id")
}

var findUserCommand = &cobra.Command{
	Use:   "find",
	Short: "finds an new user by ID or name",
	Long:  "finds an new user by ID or name",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := dig.New()

		initializers := []interface{}{
			NewConnection,
			NewClient,
		}

		for _, initFn := range initializers {
			if err := container.Provide(initFn); err != nil {
				return err
			}
		}

		// Invoke actors
		return container.Invoke(func(client user.ServiceClient) error {
			var resp *user.Users
			var err error

			if name := viper.GetString("name"); name != "" {
				resp, err = client.FindUserByName(context.Background(), &user.FindUserByNameRequest{
					Name: name,
				})
			}
			if id := viper.GetInt64("id"); id >= 0 {
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
		})
	},
}
