package jazz

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
	"time"
)

// define constraints to sql connection
const (
	maxOpenDbConn = 10 // maximum open connections at any given time
	maxDbLifetime = 5 * time.Minute
	maxIdleDbConn = 5 // how many connection can remain in the pool but idle
)

func (j *Jazz) connectToSQLDB(dbType, dsn string) (*sql.DB, error) {

	if dbType == "postgres" || dbType == "postgresql" {
		dbType = "pgx"
	}
	db, err := sql.Open(dbType, dsn)
	if err != nil {
		ErrDBNotOpen.Database = dbType
		ErrDBNotOpen.Cause = err
		return nil, ErrDBNotOpen
	}
	db.SetMaxOpenConns(maxOpenDbConn)
	db.SetMaxIdleConns(maxIdleDbConn)
	db.SetConnMaxLifetime(maxDbLifetime)

	err = db.Ping()
	if err != nil {
		ErrDBNotConnected.Database = dbType
		ErrDBNotConnected.Cause = err
		return nil, ErrDBNotConnected
	}
	log.Println("Connected to Database:", dbType)
	return db, nil
}

func (j *Jazz) connectToMongoDB(dsn string) (*mongo.Client, error) {
	log.Println("Connecting to MongoDB")

	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dsn))
	if err != nil {
		ErrDBNotOpen.Database = "mongodb"
		ErrDBNotOpen.Cause = err
		panic(ErrDBNotOpen)
	}

	err = client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		ErrDBNotConnected.Database = "mongodb"
		ErrDBNotConnected.Cause = err
		log.Fatal(ErrDBNotConnected)
	}
	log.Println("Connected to MongoDB")
	return client, nil
}

func (j *Jazz) BuildDSN() string {
	var dsn string

	switch os.Getenv("DATABASE_TYPE") {
	case "postgres", "postgresql":
		dsn = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s timezone=UTC connect_timeout=5",
			os.Getenv("DATABASE_HOST"),
			os.Getenv("DATABASE_PORT"),
			os.Getenv("DATABASE_USER"),
			os.Getenv("DATABASE_NAME"),
			os.Getenv("DATABASE_SSL_MODE"),
		)
		//left out the password field because by default is empty pass
		//if user provided a pass then append it to the dsn
		if os.Getenv("DATABASE_PASS") != "" {
			dsn = fmt.Sprintf("%s password=%s", dsn, os.Getenv("DATABASE_PASS"))
		}
	case "mongo", "mongodb":
		dsn = os.Getenv("MONGODB_URI")
	default:

	}
	return dsn
}
