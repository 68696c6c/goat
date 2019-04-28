built in migrations

default commands like make:migration, something like this:
func init() {
	RootCommand.AddCommand(goat.MakeMigrationCommand)
}

generators for application (including structure, Dockerfile and Makefile), repos, handlers, etc

cloudformation support template generator with Make targets for bootstrapping

Oauth support
