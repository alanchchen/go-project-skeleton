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
	rootCmd.AddCommand(addUserCommand)

	addUserCommand.Flags().String("name", "", "the user name")
}

var addUserCommand = &cobra.Command{
	Use:   "add",
	Short: "adds an new user",
	Long:  "adds an new user",
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
			resp, err := client.AddUser(context.Background(), &user.AddUserRequest{
				Name: viper.GetString("name"),
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
		})
	},
}
