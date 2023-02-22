package main

import (
	"errors"
	"github.com/ahmedkhaled/jazz"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"os"
)

const version = "1.0.0"

var jz jazz.Jazz

func run() {
	arg1, arg2, arg3, err := validateInput()
	if err != nil {
		exitGracefully(err)
	}
	setup(arg1)

	switch arg1 {
	case "help":
		showHelp()
	case "new":
		if arg2 == "" {
			exitGracefully(errors.New("new requires an application name"))
		}
		startNew(arg2)

	case "version":
		color.Yellow("Application version: " + version)

	case "migrate":
		if arg2 == "" {
			arg2 = "up" //by default up
		}
		err := migrate(arg2, arg3)
		if err != nil {
			exitGracefully(err)
		}

	case "create":
		if arg2 == "" {
			exitGracefully(errors.New("generate requires a subcommand: [migration|model|handler]"))
		}
		err = create(arg2, arg3)
		if err != nil {
			exitGracefully(err)
		}

	default:
		showHelp()
	}

}

func validateInput() (string, string, string, error) {
	var arg1, arg2, arg3 string

	//os.Args[0] (application name)
	if len(os.Args) == 1 {
		color.Red("Error: command required")
		showHelp()
		return "", "", "", errors.New("command required")
	}
	//check any commands other than the os.Args[0] (application name)
	if len(os.Args) >= 1 {
		arg1 = os.Args[1]
		if len(os.Args) >= 3 {
			arg2 = os.Args[2]
		}
		if len(os.Args) >= 4 {
			arg3 = os.Args[3]
		}

	}
	return arg1, arg2, arg3, nil

}
func setup(arg1 string) {
	if arg1 != "new" && arg1 != "version" && arg1 != "help" {
		err := godotenv.Load()
		if err != nil {
			exitGracefully(err)
		}

		path, err := os.Getwd()
		if err != nil {
			exitGracefully(err)
		}
		jz.DB.Type = os.Getenv("DATABASE_TYPE")

		jz.RootPath = path
	}

}
