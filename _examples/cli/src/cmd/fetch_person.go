package cmd

import (
	"encoding/json"

	"github.com/68696c6c/goat"
	"github.com/spf13/cobra"

	"github.com/68696c6c/cli/swapi"
)

func init() {
	Root.AddCommand(&cobra.Command{
		Use:   "fetch-person",
		Short: "Fetches the specified person from the Star Wars API and prints the result.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			goat.MustInit()

			id := args[0]
			println("fetching star wars person id", id)

			client, err := swapi.NewSwapi()
			if err != nil {
				return err
			}

			result, err := client.GetPerson(id)
			if err != nil {
				return err
			}

			output, err := json.MarshalIndent(result, "", "  ")
			if err != nil {
				return err
			}

			println(string(output))

			return nil
		},
	})
}
