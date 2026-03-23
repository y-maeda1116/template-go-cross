package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewHelloCommand helloコマンドを作成
func NewHelloCommand() *cobra.Command {
	var name string

	cmd := &cobra.Command{
		Use:   "hello",
		Short: "Say hello",
		Run: func(cmd *cobra.Command, args []string) {
			if name == "" {
				name = "World"
			}
			fmt.Printf("Hello, %s!\n", name)
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "Name to greet")

	return cmd
}
