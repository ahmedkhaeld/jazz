package jazz

import (
	"database/sql"
	"go.mongodb.org/mongo-driver/mongo"
)

// Settings holds all the configuration values for the jazz package
type settings struct {
	port        string // what port the application listen to
	renderer    string //what template engine to use
	sessionType string
	cookie
	databaseConfig
	redisConfig
	encryptionKey string
}

type Server struct {
	ServerName, Port string
	Secure           bool
	URL              string
}

type pathOptions struct {
	rootPath    string
	folderNames []string
}

type cookie struct {
	name     string //name of the cookie
	lifeTime string //how long the cookie lasts
	persist  string //does it persist between browser closes
	secure   string //is the cookie secure
	domain   string //what domain is the cookie associated with
}

type databaseConfig struct {
	dsn    string
	dbType string
}

type Database struct {
	Type      string
	SqlPool   *sql.DB
	MongoPool *mongo.Client
}

type redisConfig struct {
	host     string
	password string
	prefix   string
}
