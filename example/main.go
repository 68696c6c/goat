package main

import (
	"github.com/68696c6c/goat"
	"os"
	"goat/example/foo"
	"github.com/spf13/viper"
)

func main() {
	goat.Init()
	file := goat.CWD() + "/main.go"

	goat.PrintHeading("\nWelcome to the goat example app!")
	goat.PrintSuccess("Goat is a package for helping you build web apps in Go.  It is built using Gin, Viper, GORM, Logrus, and a few other open-source libraries.")
	goat.PrintSuccess("This example will walk you how Goat works.  To see how the things demonstrated here are done, take a look at the code in this file:")
	goat.PrintInfo(file)

	goat.PrintSuccess("\nThe first step to using goat is to run goat.Init()")
	goat.PrintSuccess("When you initialize Goat, it will attempt to read your project's config file.")
	goat.PrintSuccess("Doing that involves a few things.  Let's look into that.")

	// Root demo
	goat.PrintHeading("\nProject Root")
	goat.PrintSuccess("By default, Goat uses convention to find the root path of your app and read your app's config file.")
	goat.PrintSuccess("By 'root path of your app' Goat mean means 'the path to a directory containing your app's configuration, log files, and assets'.")
	goat.PrintSuccess("Of course, you can override the path completely, or manually set the path to each of those resources.")
	goat.PrintSuccess("The following is a demonstration of this.")
	demoRoot()
	goat.PrintSuccess("Even if you override the root path, Goat will always remember the path to your executable and you can access it by calling goat.ExePath()")
	ep := goat.ExePath()
	goat.PrintIndent("Executable path: " + ep)
	ed := goat.ExeDir()
	goat.PrintSuccess("If you want to get the directory holding the executable instead, you can call goat.ExeDir()")
	goat.PrintIndent("Executable directory: " + ed)

	// Config demo
	goat.PrintHeading("\nApp Configuration")
	goat.PrintSuccess("Goat expects that your config file will be named `config.yml`, in your project root directory.")
	goat.PrintSuccess("Goat will read your config using Viper.  Immediately after calling goat.Init() you can use viper to access your config values.")
	demoConfig()
	goat.PrintSuccess(`If you want to change the path to your config, you can call goat.SetConfigFile("path/to/config.ini")`)
	goat.PrintSuccess("Note that unlike most other Goat functions, this function needs to be called BEFORE you call goat.Init().")
	goat.PrintSuccess("Since goat.Init() tries to load your config, trying to override the config after it's been read will return an error.")
	goat.PrintSuccess("Notice that if we call 'goat.SetConfigFile(`/path/to/config.ini`)' now, it will have no effect.")
	goat.SetConfigFile("/path/to/config.ini")
	demoConfig()

	// CWD demo
	goat.PrintHeading("\nCurrent Working Directory")
	goat.PrintSuccess("In addition to goat.Root(), Goat provides another path helper, goat.CWD()")
	goat.PrintSuccess("This will return the path to the directory in your app where you called goat.CWD()")
	demoCWD()

	// Errors demo
	goat.PrintHeading("\nErrors")
	goat.PrintSuccess("If Goat encounters any non-fatal errors at runtime, it will add them to an internal error array.")
	goat.PrintSuccess("Goat functions that return errors will also add to this array.")
	goat.PrintSuccess("Earlier, in the config demo, we tried to set the config after the config had already been read.")
	goat.PrintSuccess("We didn't capture the error then, but we should see it if we call goat.GetErrors() now:")
	demoErrors()

	// Database demo
	goat.PrintHeading("\nDatabase")
	demoDatabase()

	os.Exit(0)
}

func demoRoot() {
	goat.PrintIndent("The default `root path` is the path to the running executable.")
	goat.PrintIndent("This path is set when goat.Init() is called and is static.  You can safely call goat.Root() from any directory in your project.")

	printRoot()

	goat.PrintIndent("Goat assumes that this root path contains all of your app's non-compiled assets.")
	goat.PrintIndent("If your executable is running from a different directory, you can manually set the project root by calling goat.SetRoot(`/some/absolute/path`).")

	goat.SetRoot("/some/absolute/path")
}

func printRoot() {
	p := goat.Root()
	goat.PrintIndent("Root path from the example main package: " + p)
	for _, p := range foo.FooRoot() {
		goat.PrintIndent(p)
	}
}

func demoCWD() {
	c := goat.CWD()
	goat.PrintIndent("CWD from the example main package: " + c)
	for _, p := range foo.FooCWD() {
		goat.PrintIndent(p)
	}
}

func demoConfig() {
	f := viper.ConfigFileUsed()
	goat.PrintIndent("Configuration read from " + f)
	v := viper.GetString("key")
	goat.PrintIndent("config.key: " + v)
}

func demoErrors() {
	errs := goat.GetErrors()
	goat.PrintIndent(goat.ErrorsToString(errs))
}

func demoDatabase() {
	println("coming soon")
}
