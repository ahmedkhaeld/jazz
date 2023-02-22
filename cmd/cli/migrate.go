package main

import (
	"fmt"
	"github.com/fatih/color"
	"os"
)

func migrate(arg2, arg3 string) error {
	dsn := buildDSN()

	// run migration scenarios
	switch arg2 {
	case "up":
		err := jz.MigrateUp(dsn)
		if err != nil {
			return err
		}

	case "down":

		if arg3 == "all" {
			err := jz.MigrateDownAll(dsn)
			if err != nil {
				return err
			}
		} else {
			err := jz.Steps(-1, dsn)
			if err != nil {
				return err
			}
		}

	case "reset":
		err := jz.MigrateDownAll(dsn)
		if err != nil {
			return err
		}
		err = jz.MigrateUp(dsn)
		if err != nil {
			return err
		}
	default:
		showHelp()
	}
	color.Green("Migration Completed!")

	return nil
}

//buildDSN builds a dsn string its syntax compatible with migrate functions based on database type
func buildDSN() string {
	dbType := jz.DB.Type
	var dsn string

	if dbType == "pgx" {
		dbType = "postgres"
	}
	if dbType == "mongo" {
		dbType = "mongodb"
	}

	if dbType == "postgres" {
		var uri string
		if os.Getenv("DATABASE_PASS") != "" {
			uri = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
				os.Getenv("DATABASE_USER"),
				os.Getenv("DATABASE_PASS"),
				os.Getenv("DATABASE_HOST"),
				os.Getenv("DATABASE_PORT"),
				os.Getenv("DATABASE_NAME"),
				os.Getenv("DATABASE_SSL_MODE"))
		} else {
			uri = fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=%s",
				os.Getenv("DATABASE_USER"),
				os.Getenv("DATABASE_HOST"),
				os.Getenv("DATABASE_PORT"),
				os.Getenv("DATABASE_NAME"),
				os.Getenv("DATABASE_SSL_MODE"))
		}
		dsn = uri

	}
	if dbType == "mysql" {
		dsn = "mysql://" + jz.BuildDSN()

	}
	if dbType == "mongodb" {
		dsn = jz.BuildDSN()
	}

	return dsn

}
