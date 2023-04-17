package cmd

import "github.com/spf13/cobra"

var Root = &cobra.Command{
	Use:   "cli",
	Short: "Root command for example CLI app.",
}
