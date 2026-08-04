package archytas

import (
	"time"

	"github.com/dgraph-io/badger/v4"
)

func (b *Badger) Close() error {
	return b.db.Close()
}

func (b *Badger) Set(key, value string) error {
	return b.SetBytes(key, []byte(value))
}

func (b *Badger) SetBytes(key string, value []byte) error {
	return b.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), value)
	})
}

func (b *Badger) SetWithTTL(key, value string, ttl time.Duration) error {
	return b.SetBytesWithTTL(key, []byte(value), ttl)
}

func (b *Badger) SetBytesWithTTL(key string, value []byte, ttl time.Duration) error {
	return b.db.Update(func(txn *badger.Txn) error {
		return txn.SetEntry(badger.NewEntry([]byte(key), value).WithTTL(ttl))
	})
}

func (b *Badger) Read(key string) ([]byte, error) {
	return b.Get(key)
}

func (b *Badger) Get(key string) ([]byte, error) {
	var copiedValue []byte
	err := b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		copiedValue, err = item.ValueCopy(nil)
		return err
	})

	return copiedValue, err
}

func (b *Badger) GetString(key string) (string, error) {
	value, err := b.Get(key)
	return string(value), err
}

func (b *Badger) Exists(key string) (bool, error) {
	err := b.db.View(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(key))
		return err
	})
	if err == badger.ErrKeyNotFound {
		return false, nil
	}
	return err == nil, err
}

func (b *Badger) Delete(key string) error {
	return b.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
}

func (b *Badger) DeletePrefix(prefix string) error {
	return b.db.DropPrefix([]byte(prefix))
}
