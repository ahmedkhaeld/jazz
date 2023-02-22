package cache

import (
	"github.com/dgraph-io/badger/v3"
	"log"
	"os"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/gomodule/redigo/redis"
)

var testRedisCache Redis
var testBadgerCache Badger

func TestMain(m *testing.M) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	redisConn := redis.Pool{
		MaxIdle:     50,
		MaxActive:   1000,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", s.Addr())
		},
	}

	testRedisCache.Conn = &redisConn
	testRedisCache.Prefix = "test-jazz"

	defer func(conn *redis.Pool) {
		err := conn.Close()
		if err != nil {

		}
	}(testRedisCache.Conn)

	//ensure to clean up badger test data after it finished
	_ = os.RemoveAll("./testdata/tmp/badger")

	// create a badger database
	if _, err := os.Stat("./testdata/tmp"); os.IsNotExist(err) {
		err := os.Mkdir("./testdata/tmp", 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = os.Mkdir("./testdata/tmp/badger", 0755)
	if err != nil {
		log.Fatal(err)
	}

	badgerConn, _ := badger.Open(badger.DefaultOptions("./testdata/tmp/badger"))
	testBadgerCache.Conn = badgerConn

	os.Exit(m.Run())
}
