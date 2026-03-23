package cli

import (
	"github.com/spf13/cobra"
)

// NewRootCommand ルートコマンドを作成
func NewRootCommand(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "myapp",
		Short: "A CLI and Desktop application",
		Long:  "MyApp is a CLI and Desktop application template built with Go and Wails.",
	}

	// サブコマンドを追加
	cmd.AddCommand(NewVersionCommand())
	cmd.AddCommand(NewHelloCommand())

	return cmd
}
