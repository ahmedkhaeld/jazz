package jazz

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
)

func (j *Jazz) MigrateUp(dsn string) error {
	m, err := migrate.New("file://"+j.RootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer func(m *migrate.Migrate) {
		err, _ := m.Close()
		if err != nil {

		}
	}(m)

	if err := m.Up(); err != nil {
		log.Println("Error running migration:", err)
		return err
	}

	return nil
}

func (j *Jazz) MigrateDownAll(dsn string) error {
	m, err := migrate.New("file://"+j.RootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer func(m *migrate.Migrate) {
		err, _ := m.Close()
		if err != nil {

		}
	}(m)

	err = m.Down()
	if err != nil {
		return err
	}

	return nil
}

func (j *Jazz) Steps(n int, dsn string) error {
	m, err := migrate.New("file://"+j.RootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer func(m *migrate.Migrate) {
		err, _ := m.Close()
		if err != nil {

		}
	}(m)

	err = m.Steps(n)
	if err != nil {
		return err
	}

	return nil
}

func (j *Jazz) MigrateForce(dsn string) error {
	m, err := migrate.New("file://"+j.RootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer func(m *migrate.Migrate) {
		err, _ := m.Close()
		if err != nil {

		}
	}(m)

	err = m.Force(-1)
	if err != nil {
		return err
	}

	return nil
}
