package main

import (
	"fmt"
	"os"

	greeter "github.com/alanchchen/go-project-skeleton/pkg/greeter/cmd"
)

func main() {
	if err := greeter.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
