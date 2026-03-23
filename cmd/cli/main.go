package main

import (
	"os"

	"github.com/y-maeda1116/template-go-cross/internal/cli"
)

func main() {
	cmd := cli.NewRootCommand("1.0.0")
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
