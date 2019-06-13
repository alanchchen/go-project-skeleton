package main

import (
	"context"
	"encoding/json"
	"fmt"

	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/alanchchen/go-project-skeleton/pkg/api/user"
	"github.com/alanchchen/go-project-skeleton/pkg/app"
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
		runner := app.NewRunner()
		if err := runner.BindCobraCommand(cmd, args...); err != nil {
			return err
		}

		initializers := []interface{}{
			NewConnection,
			NewClient,
		}

		return runner.RunCustom(func(client user.ServiceClient, cfg *viper.Viper) error {
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
