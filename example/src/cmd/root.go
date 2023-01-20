package cmd

import "github.com/spf13/cobra"

var Root = &cobra.Command{
	Use:   "example",
	Short: "Root command for example",
}
