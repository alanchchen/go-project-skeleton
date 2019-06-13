package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/cobra"

	"github.com/alanchchen/go-project-skeleton/pkg/api/user"
	"github.com/alanchchen/go-project-skeleton/pkg/app"
)

func init() {
	rootCmd.AddCommand(listUsersCommand)
}

var listUsersCommand = &cobra.Command{
	Use:   "list",
	Short: "list all users",
	Long:  "list all users",
	RunE: func(cmd *cobra.Command, args []string) error {
		runner := app.NewRunner()
		if err := runner.BindCobraCommand(cmd, args...); err != nil {
			return err
		}

		initializers := []interface{}{
			NewConnection,
			NewClient,
		}

		return runner.RunCustom(func(client user.ServiceClient) error {
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
