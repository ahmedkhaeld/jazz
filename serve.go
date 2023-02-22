package jazz

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"github.com/gomodule/redigo/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"os"
	"time"
)

// ListenAndServe start the web server
func (j *Jazz) ListenAndServe() {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%s", os.Getenv("PORT")),
		ErrorLog:          j.ErrorLog,
		Handler:           j.Routes,
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       30 * time.Second,
		ReadHeaderTimeout: time.Minute,
		WriteTimeout:      600 * time.Second,
	}

	//when we stop our application close all connections
	// close database connection explicitly if it was actually opened
	if j.DB.SqlPool != nil {
		defer func(SqlPool *sql.DB) {
			err := SqlPool.Close()
			if err != nil {

			}
		}(j.DB.SqlPool)

	}

	// close mongo connection explicitly if it was actually opened
	if j.DB.MongoPool != nil {
		defer func(Pool *mongo.Client) {
			err := Pool.Disconnect(context.Background())
			if err != nil {

			}
		}(j.DB.MongoPool)
	}

	// close redis cache connection if it was actually opened
	if redisConnection != nil {
		defer func(redisConnection *redis.Pool) {
			err := redisConnection.Close()
			if err != nil {

			}
		}(redisConnection)
	}

	// close badger cache connection if it was actually opened
	if badgerConnection != nil {
		defer func(badgerConnection *badger.DB) {
			err := badgerConnection.Close()
			if err != nil {

			}
		}(badgerConnection)
	}
	j.InfoLog.Printf("Listening on port %s", os.Getenv("PORT"))

	err := srv.ListenAndServe()
	j.ErrorLog.Fatal(err)

}
