package main

import (
	"os"

	"github.com/brevdev/cloud/cmd/bcloud/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		os.Exit(1)
	}
}
