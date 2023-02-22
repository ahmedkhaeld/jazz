package cache

import (
	"github.com/dgraph-io/badger/v3"
	"time"
)

func (b *Badger) Has(key string) (bool, error) {
	_, err := b.Get(key)
	if err != nil {
		return false, nil
	}
	return true, nil
}

func (b *Badger) Get(key string) (interface{}, error) {
	var fromCache []byte

	err := b.Conn.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}

		err = item.Value(func(val []byte) error {
			fromCache = append([]byte{}, val...)
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	decoded, err := decode(string(fromCache))
	if err != nil {
		return nil, err
	}

	item := decoded[key]

	return item, nil
}

func (b *Badger) Set(key string, value interface{}, expires ...int) error {
	entry := Entry{}

	entry[key] = value
	encoded, err := encode(entry)
	if err != nil {
		return err
	}

	if len(expires) > 0 {
		err = b.Conn.Update(func(txn *badger.Txn) error {
			e := badger.NewEntry([]byte(key), encoded).WithTTL(time.Second * time.Duration(expires[0]))
			err = txn.SetEntry(e)
			return err
		})
	} else {
		err = b.Conn.Update(func(txn *badger.Txn) error {
			e := badger.NewEntry([]byte(key), encoded)
			err = txn.SetEntry(e)
			return err
		})
	}

	return nil
}

func (b *Badger) Forget(key string) error {
	err := b.Conn.Update(func(txn *badger.Txn) error {
		err := txn.Delete([]byte(key))
		return err
	})

	return err
}

func (b *Badger) EmptyByMatch(key string) error {
	return b.emptyByMatch(key)

}

func (b *Badger) Empty() error {
	return b.emptyByMatch("")

}

func (b *Badger) emptyByMatch(str string) error {
	deleteKeys := func(keysForDelete [][]byte) error {
		if err := b.Conn.Update(func(txn *badger.Txn) error {
			for _, key := range keysForDelete {
				if err := txn.Delete(key); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			return err
		}
		return nil
	}

	collectSize := 100000

	err := b.Conn.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.AllVersions = false
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()

		keysForDelete := make([][]byte, 0, collectSize)
		keysCollected := 0

		for it.Seek([]byte(str)); it.ValidForPrefix([]byte(str)); it.Next() {
			key := it.Item().KeyCopy(nil)
			keysForDelete = append(keysForDelete, key)
			keysCollected++
			if keysCollected == collectSize {
				if err := deleteKeys(keysForDelete); err != nil {
					return err
				}
			}
		}

		if keysCollected > 0 {
			if err := deleteKeys(keysForDelete); err != nil {
				return err
			}
		}

		return nil
	})

	return err
}
