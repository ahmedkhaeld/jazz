package jazz

import (
	"github.com/ahmedkhaeld/jazz/cache"
	"github.com/ahmedkhaeld/jazz/mailer"
	"github.com/ahmedkhaeld/jazz/render"
	"github.com/ahmedkhaeld/jazz/session"
	"github.com/alexedwards/scs/v2"
	"github.com/dgraph-io/badger/v3"
	"github.com/go-chi/chi/v5"
	"github.com/gomodule/redigo/redis"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"log"
	"os"
	"strconv"
	"strings"
)

const version = "1.0.0"
const randomString = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPGRSTUVWXYZ0987654321_+"

// Jazz is the overall type for the jazz package.
//
// define the requirements of using a jazz project
type Jazz struct {
	AppName       string              // application name
	Debug         bool                // debug mode
	Version       string              // version of the application
	ErrorLog      *log.Logger         // error logger
	InfoLog       *log.Logger         // info logger
	RootPath      string              // root path of the application
	Routes        *chi.Mux            // router
	Render        *render.Render      // render engine
	Session       *scs.SessionManager // session manager accessible from the top level
	DB            Database            // database
	DSN           string              // database connection string
	settings                          // application settings
	EncryptionKey string              // encryption key
	Cache         cache.Cache         // cache
	Scheduler     *cron.Cron          // scheduler
	Mailer        mailer.Mail         // mailer
	Server        Server              // server
}

// New instantiate a Jazz app with its requirements
// that we expect to be available at the top level of the application
func (j *Jazz) New(rootPath string) error {

	err := j.initDefaultPaths(rootPath)
	if err != nil {
		return err
	}
	// make sure there is .env file
	err = j.CreateDotEnvIfNotExists(rootPath)
	if err != nil {
		return err
	}

	//read .env
	err = godotenv.Load(rootPath + "/.env")
	if err != nil {
		ErrPathNotFound.Cause = err
		return ErrPathNotFound
	}

	//populate settings fields with .env file values
	j.settings = j.envSettings()
	InfoL, ErrL := j.createLoggers(os.Stdout)
	j.InfoLog, j.ErrorLog = InfoL, ErrL
	j.AppName = os.Getenv("APP_NAME")
	j.Debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))
	j.RootPath = rootPath
	j.Version = version
	j.DB = j.db()
	j.Cache = j.cache()
	j.DSN = j.BuildDSN()
	j.Mailer = j.mailer()
	j.Routes = j.mux().(*chi.Mux) //set the available default routes
	j.Server = j.server()
	j.Session = j.session()
	j.EncryptionKey = j.settings.encryptionKey
	j.Render = j.render()

	go j.Mailer.ListenForMail()

	return nil
}

// defaultPaths provide the application with predefined directories if they not exists
func (j *Jazz) initDefaultPaths(root string) error {
	opt := pathOptions{
		folderNames: []string{
			"handlers", "migrations", "mail", "views", "data", "public", "tmp", "logs", "middleware",
		},
	}
	for _, path := range opt.folderNames {
		//create folder if it doesn't exist
		err := j.CreateDirIfNotExist(root + "/" + path)
		if err != nil {
			return err
		}
	}
	return nil
}
func (j *Jazz) envSettings() settings {
	return settings{
		port:     os.Getenv("PORT"),
		renderer: os.Getenv("RENDERER"),
		cookie: cookie{
			name:     os.Getenv("COOKIE_NAME"),
			lifeTime: os.Getenv("COOKIE_LIFETIME"),
			persist:  os.Getenv("COOKIE_PERSISTS"),
			secure:   os.Getenv("COOKIE_SECURE"),
			domain:   os.Getenv("COOKIE_DOMAIN"),
		},
		sessionType: os.Getenv("SESSION_TYPE"),
		databaseConfig: databaseConfig{
			dsn:    j.BuildDSN(),
			dbType: os.Getenv("DATABASE_TYPE"),
		},
		encryptionKey: os.Getenv("KEY"),
		redisConfig: redisConfig{
			host:     os.Getenv("REDIS_HOST"),
			password: os.Getenv("REDIS_PASSWORD"),
			prefix:   os.Getenv("REDIS_PREFIX"),
		},
	}

}
func (j *Jazz) server() Server {
	secure := true
	if strings.ToLower("SECURE") == "false" {
		secure = false
	}
	return Server{
		ServerName: os.Getenv("SERVER_NAME"),
		Port:       os.Getenv("PORT"),
		Secure:     secure,
		URL:        os.Getenv("APP_URL"),
	}
}
func (j *Jazz) render() *render.Render {
	rn := render.Render{
		Renderer: j.settings.renderer,
		RootPath: j.RootPath,
		Debug:    j.Debug,
		Port:     j.settings.port,
		Session:  j.Session,
	}
	return &rn

}
func (j *Jazz) session() *scs.SessionManager {
	//set the session settings
	s := session.Session{
		Cookie: session.Cookie{
			LifeTime: j.settings.cookie.lifeTime,
			Name:     j.settings.cookie.name,
			Persist:  j.settings.cookie.persist,
			Domain:   j.settings.cookie.domain,
			Secure:   j.settings.cookie.secure,
		},
		SessionType: j.settings.sessionType,
	}

	switch j.settings.sessionType {
	case "redis":
		s.RedisPool = myRedisCache.Conn
	case "mysql", "postgres", "mariadb", "postgresql":
		s.DBPool = j.DB.SqlPool
	}

	//create the session manager
	return s.New()

}
func (j *Jazz) mailer() mailer.Mail {
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	m := mailer.Mail{
		Domain:      os.Getenv("MAIL_DOMAIN"),
		Templates:   j.RootPath + "/mail",
		Host:        os.Getenv("SMTP_HOST"),
		Port:        port,
		Username:    os.Getenv("SMTP_USERNAME"),
		Password:    os.Getenv("SMTP_PASSWORD"),
		Encryption:  os.Getenv("SMTP_ENCRYPTION"),
		FromName:    os.Getenv("FROM_NAME"),
		FromAddress: os.Getenv("FROM_ADDRESS"),
		Jobs:        make(chan mailer.Message, 20),
		Results:     make(chan mailer.Result, 20),
		API:         os.Getenv("MAILER_API"),
		APIKey:      os.Getenv("MAILER_KEY"),
		APIUrl:      os.Getenv("MAILER_URL"),
	}
	return m
}

func (j *Jazz) db() Database {
	dbType := os.Getenv("DATABASE_TYPE")
	if dbType == "postgres" || dbType == "mysql" || dbType == "mariadb" {
		db, err := j.connectToSQLDB(dbType, j.BuildDSN())
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		return Database{
			Type:    dbType,
			SqlPool: db,
		}
	}
	if dbType == "mongodb" || dbType == "mongo" {
		db, err := j.connectToMongoDB(j.BuildDSN())
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		return Database{
			Type:      dbType,
			MongoPool: db,
		}
	}
	if dbType == "" {
		log.Println("No database type set")
	}
	return Database{}
}

var (
	myRedisCache     *cache.Redis
	myBadgerCache    *cache.Badger
	redisConnection  *redis.Pool
	badgerConnection *badger.DB
)

// Cache decide to or not to use a cache whether it's redis or badger
func (j *Jazz) cache() cache.Cache {
	cacheSet := os.Getenv("CACHE")
	sessionType := os.Getenv("SESSION_TYPE")
	if cacheSet == "redis" || sessionType == "redis" {
		myRedisCache = j.connectToRedis()
		redisConnection = myRedisCache.Conn

		log.Println("redis cache is set")
		return myRedisCache
	}

	scheduler := cron.New()
	j.Scheduler = scheduler
	if cacheSet == "badger" {
		myBadgerCache = j.connectToBadger()
		badgerConnection = myBadgerCache.Conn
		//schedule garbage collection once a day for housekeeping on badger
		_, err := j.Scheduler.AddFunc("@daily", func() {
			_ = myBadgerCache.Conn.RunValueLogGC(0.7)
		})
		if err != nil {
			log.Fatal(err)
		}
		log.Print("badger cache is set")
		return myBadgerCache

	}
	return nil

}
