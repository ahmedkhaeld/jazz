package main

import (
	"github.com/fatih/color"
	"os"
)

func exitGracefully(err error, msg ...string) {
	message := ""
	if len(msg) > 0 {
		message = msg[0]
	}

	if err != nil {
		color.Red("Error: %v\n", err)
	}

	if len(message) > 0 {
		color.Yellow(message)
	} else {
		color.Green("Finished!")
	}

	os.Exit(0)
}

func showHelp() {
	color.Yellow(`Available commands:

	help                  - show the help commands
	version               - print application version
	new <appName>           - create a new application
	migrate               - runs all up migrations that have not been run previously
	migrate down          - reverses the most recent migration
	migrate reset         - runs all down migrations in reverse order, and then all up migrations
	create migration <name> - creates two new up and down migrations in the migrations folder
	create auth             - creates and runs migrations for authentication tables, and creates models and middleware
	create handler <name>   - creates a stub handler in the handlers directory
	create model <name>     - creates a new model in the data directory
	create mail <name>      - creates two starter mail templates in the mail directory
	create session          - creates a table in the database as a session store
	
	`)
}
