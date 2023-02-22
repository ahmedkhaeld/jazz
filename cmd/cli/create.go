package main

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
	"os"
	"strings"
	"time"
)

func create(arg2, arg3 string) error {

	switch arg2 {
	case "migration":
		migrations(arg3)

	case "auth":
		auth()

	case "handler":
		handler(arg3)

	case "model":
		model(arg3)
	case "session":
		err := sessionsTable()
		if err != nil {
			exitGracefully(err)
		}
	case "key":
		rnd := jz.RandomString(32)
		color.Yellow("32 character encryption key %s", rnd)

	case "mail":
		mail(arg3)
	}

	return nil
}

func migrations(nameofMigration string) {
	var dbType string
	dbType = jz.DB.Type
	if nameofMigration == "" {
		exitGracefully(errors.New("you must give the migration a name"))
	}

	fileName := fmt.Sprintf("%d_%s", time.Now().UnixMicro(), nameofMigration)

	if dbType == "mongo" || dbType == "mongodb" {
		dbType = "mongo"
		upFile := jz.RootPath + "/migrations/" + fileName + "." + dbType + ".up.json"
		downFile := jz.RootPath + "/migrations/" + fileName + "." + dbType + ".down.json"
		err := copyFileFromTemplate("templates/migrations/migration."+dbType+".up.json", upFile)
		if err != nil {
			exitGracefully(err)
		}
		err = copyFileFromTemplate("templates/migrations/migration."+dbType+".down.json", downFile)
		if err != nil {
			exitGracefully(err)
		}
		color.Yellow("%s.%s migration created!", fileName, dbType)
	} else {
		upFile := jz.RootPath + "/migrations/" + fileName + "." + dbType + ".up.sql"
		downFile := jz.RootPath + "/migrations/" + fileName + "." + dbType + ".down.sql"
		err := copyFileFromTemplate("templates/migrations/migration."+dbType+".up.sql", upFile)
		if err != nil {
			exitGracefully(err)
		}
		err = copyFileFromTemplate("templates/migrations/migration."+dbType+".down.sql", downFile)
		if err != nil {
			exitGracefully(err)
		}
		color.Yellow("%s.%s migration created!", fileName, dbType)
	}

}

func auth() {

	//1-copy data
	err := copyFileFromTemplate("templates/data/user.go.txt", jz.RootPath+"/data/user.go")
	if err != nil {
		exitGracefully(err)
	}
	err = copyFileFromTemplate("templates/data/token.go.txt", jz.RootPath+"/data/token.go")
	if err != nil {
		exitGracefully(err)
	}
	err = copyFileFromTemplate("templates/data/remember_token.go.txt", jz.RootPath+"/data/remember_token.go")
	if err != nil {
		exitGracefully(err)
	}

	//2-copy handlers
	err = copyFileFromTemplate("templates/handlers/auth.go.txt.txt", jz.RootPath+"/handlers/auth.go.txt")
	if err != nil {
		exitGracefully(err)
	}

	//3-copy mailer
	err = copyFileFromTemplate("templates/mailer/password-reset.html.tmpl", jz.RootPath+"/mail/password-reset.html.tmpl")
	if err != nil {
		exitGracefully(err)
	}
	err = copyFileFromTemplate("templates/mailer/password-reset.plain.tmpl", jz.RootPath+"/mail/password-reset.plain.tmpl")
	if err != nil {
		exitGracefully(err)
	}

	//4-copy middleware
	err = copyFileFromTemplate("templates/middleware/auth.go.txt.txt", jz.RootPath+"/middleware/auth.go.txt")
	if err != nil {
		exitGracefully(err)
	}
	err = copyFileFromTemplate("templates/middleware/auth-token.go.txt.txt", jz.RootPath+"/middleware/auth-token.go.txt")
	if err != nil {
		exitGracefully(err)
	}

	err = copyFileFromTemplate("templates/middleware/remember.go.txt.txt", jz.RootPath+"/middleware/remember.go.txt")
	if err != nil {
		exitGracefully(err)
	}

	//generate migrations up&down for auth tables
	dbType := jz.DB.Type
	fileName := fmt.Sprintf("%d_create_auth_tables", time.Now().UnixMicro())
	upFile := jz.RootPath + "/migrations/" + fileName + "." + dbType + ".up.sql"
	downFile := jz.RootPath + "/migrations/" + fileName + "." + dbType + ".down.sql"

	//5-copy migrations
	err = copyFileFromTemplate("templates/migrations/auth_tables."+dbType+".up.sql", upFile)
	if err != nil {
		exitGracefully(err)
	}

	err = copyFileFromTemplate("templates/migrations/auth_tables."+dbType+".down.sql", downFile)
	if err != nil {
		exitGracefully(err)
	}

	// run migrations
	err = migrate("up", "")
	if err != nil {
		exitGracefully(err)
	}

	//6-copy views
	err = copyFileFromTemplate("templates/views/login.jet", jz.RootPath+"/views/login.jet")
	if err != nil {
		exitGracefully(err)
	}

	err = copyFileFromTemplate("templates/views/forgot.jet", jz.RootPath+"/views/forgot.jet")
	if err != nil {
		exitGracefully(err)
	}

	err = copyFileFromTemplate("templates/views/reset-password.jet", jz.RootPath+"/views/reset-password.jet")
	if err != nil {
		exitGracefully(err)
	}
	color.Yellow(" -user, tokens, and remember_tokens migrations created and executed")
	color.Yellow(" -user, token models created")
	color.Yellow(" -auth middleware create")
	color.Yellow("")
	color.Yellow("You should include User and Token types in Models type in data/models.go," +
		"And add the appropriate middleware to your routes!")

}

// Handler creates a stub handler for the end user
func handler(arg3 string) {
	if arg3 == "" {
		exitGracefully(errors.New("you must give the handler a name"))
	}
	fileName := jz.RootPath + "/handlers/" + strings.ToLower(arg3) + ".go"
	if fileExists(fileName) {
		exitGracefully(errors.New(fileName + " already exists!"))
	}

	//read data of src files from templateFS
	data, err := templateFS.ReadFile("templates/handlers/handler.go.txt")
	if err != nil {
		exitGracefully(err)
	}

	//cast data content into string
	handler := string(data)
	//search for $HANDLER$ to replace its value into CamelCase
	handler = strings.ReplaceAll(handler, "$HANDLER$", strcase.ToCamel(arg3))

	err = os.WriteFile(fileName, []byte(handler), 0644)
	if err != nil {
		exitGracefully(err)
	}
	color.Yellow("%s handler created!", arg3)

}

func model(arg3 string) {
	if arg3 == "" {
		exitGracefully(errors.New("you must give the model a name"))
	}

	data, err := templateFS.ReadFile("templates/data/model.go.txt")
	if err != nil {
		exitGracefully(err)
	}

	model := string(data)

	plur := pluralize.NewClient()

	var modelName = arg3
	var tableName = arg3

	if plur.IsPlural(arg3) {
		modelName = plur.Singular(arg3)
		tableName = strings.ToLower(tableName)
	} else {
		tableName = strings.ToLower(plur.Plural(arg3))
	}

	fileName := jz.RootPath + "/data/" + strings.ToLower(modelName) + ".go"
	if fileExists(fileName) {
		exitGracefully(errors.New(fileName + " already exists!"))
	}

	model = strings.ReplaceAll(model, "$MODELNAME$", strcase.ToCamel(modelName))
	model = strings.ReplaceAll(model, "$TABLENAME$", tableName)

	err = copyDataToFile([]byte(model), fileName)
	if err != nil {
		exitGracefully(err)
	}
	color.Yellow("%s model created!", arg3)
}

// SessionsTable creates a table named sessions in the database store, holds the session values
func sessionsTable() error {
	//figure out which database used specified in .env
	dbType := jz.DB.Type

	//cast db type name to only postgres or mysql
	if dbType == "mariadb" {
		dbType = "mysql"
	}
	if dbType == "postgresql" {
		dbType = "postgres"
	}
	fileName := fmt.Sprintf("%d_create_sessions_table", time.Now().UnixMicro())

	upFile := jz.RootPath + "/migrations/" + fileName + "." + dbType + ".up.sql"
	downFile := jz.RootPath + "/migrations/" + fileName + "." + dbType + ".down.sql"

	err := copyFileFromTemplate("templates/migrations/"+dbType+"_session.sql", upFile)
	if err != nil {
		exitGracefully(err)
	}

	err = copyDataToFile([]byte("drop table sessions;"), downFile)
	if err != nil {
		exitGracefully(err)
	}

	err = migrate("up", "")
	if err != nil {
		exitGracefully(err)
	}

	color.Yellow("sessions table created!")
	return nil

}

func mail(arg3 string) {
	if arg3 == "" {
		exitGracefully(errors.New("you must give the mail template a name"))
	}
	htmlMail := jz.RootPath + "/mail/" + strings.ToLower(arg3) + ".html.tmpl"
	plainMail := jz.RootPath + "/mail/" + strings.ToLower(arg3) + ".plain.tmpl"

	err := copyFileFromTemplate("templates/mailer/mail.html.tmpl", htmlMail)
	if err != nil {
		exitGracefully(err)
	}

	err = copyFileFromTemplate("templates/mailer/mail.plain.tmpl", plainMail)
	if err != nil {
		exitGracefully(err)
	}
}
