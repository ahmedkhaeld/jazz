package jazz

import (
	"github.com/ahmedkhaeld/jazz/cache"
	"github.com/dgraph-io/badger/v3"
	"github.com/gomodule/redigo/redis"
	"time"
)

func (j *Jazz) connectToBadger() *cache.Badger {
	cacheClient := cache.Badger{
		Conn: j.openBadgerConn(),
	}
	return &cacheClient
}

func (j *Jazz) connectToRedis() *cache.Redis {
	return &cache.Redis{
		Conn:   j.openRedisConn(),
		Prefix: j.settings.redisConfig.prefix,
	}
}

func (j *Jazz) openBadgerConn() *badger.DB {
	db, err := badger.Open(badger.DefaultOptions(j.RootPath + "/tmp/badger"))
	if err != nil {
		return nil
	}
	return db
}

func (j *Jazz) openRedisConn() *redis.Pool {
	return &redis.Pool{
		MaxIdle:      50,
		MaxActive:    10000,
		IdleTimeout:  240 * time.Second,
		Dial:         j.dialRedis,
		TestOnBorrow: ping,
	}
}

// dialRedis is a helper function to connect to redis
func (j *Jazz) dialRedis() (redis.Conn, error) {
	c, err := redis.Dial("tcp",
		j.settings.redisConfig.host,
		redis.DialPassword(j.settings.redisConfig.password))
	if err != nil {
		ErrDBDial.Database = "redis"
		ErrDBDial.Cause = err
		return nil, ErrDBDial
	}
	return c, nil
}
func ping(conn redis.Conn, t time.Time) error {
	_, err := conn.Do("PING")
	if err != nil {
		ErrDBPing.Database = "redis"
		ErrDBPing.Cause = err
		return ErrDBPing
	}
	return nil
}
