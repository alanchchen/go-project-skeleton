package main

import (
	"fmt"
	"os"

	user "github.com/alanchchen/go-project-skeleton/pkg/user/cmd"
)

func main() {
	if err := user.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
