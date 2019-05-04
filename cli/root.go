package cli

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Goat = &cobra.Command{
	Use:   "goat",
	Short: "Root command for Goat CLI",
}

func init() {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	viper.SetDefault("author", "Aaron Hill <68696c6c@gmail.com>")
	viper.SetDefault("license", "MIT")
}
