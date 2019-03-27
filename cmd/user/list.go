package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/cobra"
	"go.uber.org/dig"

	"github.com/alanchchen/go-project-skeleton/pkg/api/user"
)

func init() {
	rootCmd.AddCommand(listUsersCommand)
}

var listUsersCommand = &cobra.Command{
	Use:   "list",
	Short: "list all users",
	Long:  "list all users",
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
		})
	},
}
