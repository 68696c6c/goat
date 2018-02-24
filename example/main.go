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

	// Logging demo
	goat.PrintHeading("\nLogging")
	goat.PrintSuccess("Goat uses Logrus to handling logging.")
	demoLogging()

	// Database demo
	goat.PrintHeading("\nDatabase")
	goat.PrintSuccess("Goat uses GORM for interacting with databases.")
	demoDatabase()

	// Migrations demo
	goat.PrintHeading("\nMigrations")
	demoMigrations()

	goat.PrintHeading("\nFuture Development")
	goat.PrintSuccess("Although Goat is inspired by frameworks like Rails and Laravel and seeks to provide a similarly streamlined development process, it is also inspired by the features and style of Go itself.  Goat aims to provide maximum functionality and ease of use in a small footprint with no magic.")

	goat.PrintSuccess("\nGoat is still in the early stages of development.  Before Goat can be considered feature-complete and ready for use, features like migrations, logging, request binding, and command line tools (similar to Artisan or Rake) need to be completed.")
	goat.PrintSuccess("The Goat modules also need to be unit-tested and re-written to use interfaces, ensuring that each module adheres to a consistent pattern.")
	goat.PrintSuccess("Goat might eventually use go:generate to flesh out certain features, like adding .String() functions to some of its structs, or to add functionality to apps using Goat.")

	goat.PrintSuccess("\nGoat believes that if you build something, but don't document it, you haven't actually built anything at all.  Therefore, Goat needs to be thoroughly documented and provide features for helping it's users document their code as well.  Good documentation should include in-context examples, development guides, and API documentation.")
	goat.PrintSuccess("Goat will probably use open source libraries like GoSwagger to help meet these goals.")

	println()

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

func demoLogging() {
	goat.PrintSuccess("Goat provides a default system log that will write to [project root]/logs/sys.log")
	goat.PrintSuccess("You can change the destination of this log by adding a logs.sys entry to your app's config.yml.  Future versions of Goat will provide more configuration options.")
	goat.PrintIndent("logs:\nsys: /path/to/logs/sys.log")
	goat.PrintSuccess("You can get a pointer to the system logger by calling goat.NewLogger()")

	_ = goat.NewLogger()

	goat.PrintSuccess("You can also create new logger instances by calling goat.NewCustomLogger(logName)")
	goat.PrintSuccess("If the log name you specify exists under the 'logs' key in your config, Goat will create the log file at that path.  Otherwise it will fall back to the project root.")
	goat.PrintWarning("In the near future, Goat will allow you to set a custom log root path in the same way you can set a custom config path using 'goat.SetConfigFile(`/path/to/config.yml`)'")

	goat.SetRoot(goat.ExeDir())
	_ = goat.NewCustomLogger("log_name")

	goat.PrintSuccess("\nGoat believes good logging is a critical part of any serious framework and although logging support is currently rather primitive, future developments will focus on improving logging.")
	goat.PrintWarning("More logging examples coming soon.")
}

func demoDatabase() {
	goat.PrintSuccess("Once Goat has read your config file, you can call goat.NewDB() to try to get a db connection using config values under `db`.")
	goat.PrintSuccess("Goat uses a DBConfig struct to represent a set of database connection info.  You can get the values loaded from your config file by calling goat.GetDefaultDBConfig().")

	d := goat.GetDefaultDBConfig()
	goat.PrintIndent(d.String())

	goat.PrintSuccess("By default, Goat will panic if it fails to connect to your database.  Goat treats configuration issues like this as critical errors that will prevent your app from functioning.")
	goat.PrintSuccess("Of course, whether or not that is actually the case is up to you.  Goat aims to have strict conventions, but also provide an easy way of working around them if you need.")
	goat.PrintSuccess("We can tell Goat to just return an error instead of panicking by calling goat.SetDBPanicMode(false) before calling any other Goat database functions.")
	goat.PrintSuccess("Since the example credentials in our config file are just examples, goat will panic if we try and use them to connect.")
	goat.SetDBPanicMode(false)

	_, err1 := goat.NewDB()

	goat.PrintSuccess("You can also get a new arbitrary connection using a goat.DBConfig struct and goat.NewCustomDB(c)")
	goat.PrintSuccess("We can manually create a new DBConfig using whatever values we like.")
	c := goat.DBConfig{
		Host:     "127.0.0.1",
		Port:     3306,
		Database: "db_name",
		Username: "user",
		Password: "password",
		Debug:    true,
	}
	goat.PrintIndent(c.String())

	goat.PrintSuccess("Connecting to a database using goat.NewCustomDB(c) works exactly like using goat.NewDB().")
	goat.PrintSuccess("In fact, goat.NewDB() calls goat.NewCustomDB(c) under the hood, so goat.SetDBPanicMode() has the same affect on both functions.")

	_, err2 := goat.NewCustomDB(c)

	goat.PrintSuccess("Since we called so goat.SetDBPanicMode(false) earlier, both of our database connections should have returned errors:")

	goat.PrintIndent(err1.Error())
	goat.PrintIndent(err2.Error())

	goat.PrintSuccess("The database connections returned by goat.NewDB() and goat.NewCustomDB(c) are pointers to a gorm.DB instance.  For more information how how to use them, consult the GORM documentation.")
	goat.PrintSuccess("Future developments will add the option to use a different ORM and include more in-depth examples on how to build database driven apps.")
}

func demoMigrations() {
	goat.PrintSuccess("Goat believes that a framework that doesn't include database migrations is not a complete framework and that migrations should be explicit, incremental, reversible, and divorced from your actual data models.")
	goat.PrintSuccess("GORM supports auto-migrations, using your model struct definitions to generate tables.  This feature can be useful for quick-and-dirty development, but Goat seeks to provide another way of handling migrations that will feel more familiar to users of other MVC frameworks like Laravel and Rails.")
	goat.PrintSuccess("To that end, Goat will use a standalone library called Goose for database migrations.  Goose is a separate project and you will need to import it yourself if you want to use it.  In the near future, Goat will come with Goose ready to use out of the box.")
	goat.PrintSuccess("Like Goat, Goose is built using open-source libraries like GORM and is in the very early stages of development.  Also like Goat, Goose aims to be simple to use and convention driven without getting in your way.")
	goat.PrintSuccess("Currently Goose only supports MySQL, but future developments will expand this to include all database drivers supported by Go.")
	goat.PrintSuccess("Goose provides a Schema struct with Install, Up, Drop, Reset, and CreateMigration functions.")
	goat.PrintWarning("Examples of how to use Goose will be added here soon, but for now you can look at the Goose code here: github.com/68696c6c/goose")
	// @TODO demo migrations
}