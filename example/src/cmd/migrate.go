package cmd

import (
	"github.com/68696c6c/goat"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"

	_ "github.com/68696c6c/example/db/migrations"
)

func init() {
	Root.AddCommand(&cobra.Command{
		Use:   "migrate",
		Short: "goose migrations (go run main.go migrate up)",
		Args:  cobra.MinimumNArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			goat.Init()

			db, err := goat.GetMainDB()
			if err != nil {
				goat.ExitError(errors.Wrap(err, "error initializing migration connection"))
			}

			if err := goose.SetDialect("mysql"); err != nil {
				goat.ExitError(errors.Wrap(err, "error initializing goose"))
			}

			var arguments []string
			if len(args) > 1 {
				arguments = args[1:]
			}

			sqlDb, err := db.DB()
			if err != nil {
				goat.ExitError(err)
			}

			if err := goose.Run(args[0], sqlDb, ".", arguments...); err != nil {
				goat.ExitError(err)
			}

			goat.ExitSuccess()
		},

		Example: `
Usage: app migrate [OPTIONS] COMMAND

Drivers:
postgres
mysql
sqlite3
redshift

Commands:
up                   Migrate the DB to the most recent version available
up-to VERSION        Migrate the DB to a specific VERSION
down                 Roll back the version by 1
down-to VERSION      Roll back to a specific VERSION
redo                 Re-run the latest migration
status               Dump the migration status for the current DB
version              Print the current version of the database
create NAME [sql|go] Creates new migration file with the current timestamp

Examples:
app migrate status
app migrate create init sql
app migrate create add_some_column sql
app migrate create fetch_user_data go
app migrate up

app migrate status`,
	})
}
