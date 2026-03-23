package main

import (
	"os"

	"github.com/user/repo/internal/cli"
)

func main() {
	cmd := cli.NewRootCommand("1.0.0")
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
