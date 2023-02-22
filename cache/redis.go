package cache

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

func (c *Redis) Has(key string) (bool, error) {
	key = fmt.Sprintf("%s:%s", c.Prefix, key)
	conn := c.Conn.Get()
	defer func(conn redis.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	ok, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false, err
	}

	return ok, nil
}

// Get something from cache based on a compound key bind with prefix specified only for the application
func (c *Redis) Get(key string) (interface{}, error) {
	key = fmt.Sprintf("%s:%s", c.Prefix, key)
	conn := c.Conn.Get()
	defer func(conn redis.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	cacheEntry, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}

	decoded, err := decode(string(cacheEntry))
	if err != nil {
		return nil, err
	}

	item := decoded[key]

	return item, nil
}

func (c *Redis) Set(key string, value interface{}, expires ...int) error {
	key = fmt.Sprintf("%s:%s", c.Prefix, key)
	conn := c.Conn.Get()
	defer func(conn redis.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	entry := Entry{}
	entry[key] = value
	encoded, err := encode(entry)
	if err != nil {
		return err
	}

	if len(expires) > 0 {
		_, err := conn.Do("SETEX", key, expires[0], string(encoded))
		if err != nil {
			return err
		}
	} else {
		_, err := conn.Do("SET", key, string(encoded))
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Redis) Forget(key string) error {
	key = fmt.Sprintf("%s:%s", c.Prefix, key)
	conn := c.Conn.Get()
	defer func(conn redis.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	_, err := conn.Do("DEL", key)
	if err != nil {
		return err
	}

	return nil
}

func (c *Redis) EmptyByMatch(key string) error {
	key = fmt.Sprintf("%s:%s", c.Prefix, key)
	conn := c.Conn.Get()
	defer conn.Close()

	keys, err := c.getKeys(key)
	if err != nil {
		return err
	}

	for _, x := range keys {
		_, err := conn.Do("DEL", x)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Redis) Empty() error {
	key := fmt.Sprintf("%s:", c.Prefix)
	conn := c.Conn.Get()
	defer conn.Close()

	keys, err := c.getKeys(key)
	if err != nil {
		return err
	}

	for _, k := range keys {
		_, err := conn.Do("DEL", k)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Redis) getKeys(pattern string) ([]string, error) {
	conn := c.Conn.Get()
	defer func(conn redis.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	iter := 0
	var keys []string

	for {
		arr, err := redis.Values(conn.Do("SCAN", iter, "MATCH", fmt.Sprintf("%s*", pattern)))
		if err != nil {
			return keys, err
		}

		iter, _ = redis.Int(arr[0], nil)
		k, _ := redis.Strings(arr[1], nil)
		keys = append(keys, k...)

		if iter == 0 {
			break
		}
	}

	return keys, nil
}
