package cmd

import "github.com/spf13/cobra"

var Root = &cobra.Command{
	Use:   "web",
	Short: "Root command for the example web app.",
}
