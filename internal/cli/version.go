package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewVersionCommand バージョンコマンドを作成
func NewVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Prints version number",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("MyApp v1.0.0")
		},
	}
}
