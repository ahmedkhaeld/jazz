package cache

import (
	"bytes"
	"encoding/gob"
	"github.com/dgraph-io/badger/v3"
	"github.com/gomodule/redigo/redis"
)

type Cache interface {
	Has(string) (bool, error)              //does my cache has some string
	Get(string) (interface{}, error)       //get something from the cache
	Set(string, interface{}, ...int) error //to store something in the cache
	Forget(string) error                   //get out of the cache
	EmptyByMatch(string) error             //forget everything in the cache by a pattern
	Empty() error                          //forget all
}

type Redis struct {
	Conn   *redis.Pool
	Prefix string
}

type Badger struct {
	Conn   *badger.DB
	Prefix string
}

//Entry define data serialization format of key string associated with any kind of data
type Entry map[string]interface{}

//encode serialize an item
func encode(item Entry) ([]byte, error) {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(item)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

//decode deserialize an item
func decode(str string) (Entry, error) {
	item := Entry{}
	b := bytes.Buffer{}
	b.Write([]byte(str))
	d := gob.NewDecoder(&b)
	err := d.Decode(&item)
	if err != nil {
		return nil, err
	}
	return item, nil
}
