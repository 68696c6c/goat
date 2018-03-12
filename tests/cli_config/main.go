package main

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/68696c6c/goat"
	"github.com/gin-gonic/gin"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var configFile string

var cmd = &cobra.Command{
	Use: "test",
	Run: func(cmd *cobra.Command, args []string) {
		port := fmt.Sprintf(":%s", viper.GetString("listen-port"))
		env := viper.GetString("env")

		router := goat.NewRouter(env, func(r *goat.Router) {
			r.Engine.GET("/ping", func(c *gin.Context) {
				goat.RespondMessage(c, "pong")
			})
		})

		err := router.Run(port)
		if err != nil {
			panic(fmt.Sprintf("Error starting server: %s\n", err))
		}
	},
}

func init() {
	cobra.OnInitialize(initGoat)
	viper.SetDefault("configFile", "./config.yml")
	cmd.PersistentFlags().StringVar(&configFile, "config", "./config.yml", "config file (default is ./config.yml)")
	viper.BindPFlag("projectBase", cmd.PersistentFlags().Lookup("projectBase"))
	viper.BindPFlag("useViper", cmd.PersistentFlags().Lookup("viper"))
	viper.SetDefault("author", "Aaron Hill <aaron.hill@cypcorp.com>")
	viper.SetDefault("license", "Copyright 2017 CYP")
}

func initGoat() {
	goat.SetConfigFilePath(configFile)
	goat.Init()
}
