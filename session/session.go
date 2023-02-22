package session

import (
	"database/sql"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gomodule/redigo/redis"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Session is the session manager type for jazz type to use it as a dependency
type Session struct {
	Cookie
	SessionType string
	DBPool      *sql.DB
	RedisPool   *redis.Pool
}
type Cookie struct {
	LifeTime string
	Secure   string
	Persist  string
	Name     string
	Domain   string
}

// New create a new session manager
// it will create a session manager based on the session type in the .env file
// with server side store [redis, mysql, postgres] or client side store [cookie]
func (s *Session) New() *scs.SessionManager {
	var persist, secure bool

	//cookie configurations value casting from .env file string to bool
	if strings.ToLower(s.Cookie.Persist) == "true" {
		persist = true
	}
	if strings.ToLower(s.Secure) == "true" {
		secure = true
	}
	//cookie configurations value casting from .env file string to int
	minutes, err := strconv.Atoi(s.Cookie.LifeTime)
	if err != nil {
		minutes = 60
	}
	//create session
	session := scs.New()
	session.Lifetime = time.Duration(minutes) * time.Minute
	session.Cookie.Name = s.Cookie.Name
	session.Cookie.Persist = persist
	session.Cookie.Secure = secure
	session.Cookie.Domain = s.Cookie.Domain
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.HttpOnly = true

	// which session store? [server-side] redis, mysql, postgres or cookie[client-side]
	switch strings.ToLower(s.SessionType) {
	case "redis":
		session.Store = redisstore.New(s.RedisPool)
	case "mysql", "mariadb":
		session.Store = mysqlstore.New(s.DBPool)

	case "postgres", "postgresql":
		session.Store = postgresstore.New(s.DBPool)
	default:
		// cookie
	}

	return session
}
