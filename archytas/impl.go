package cache

import (
	"github.com/dgraph-io/badger/v3"
	uuid2 "github.com/google/uuid"
	"path/filepath"
	"time"
)

type Client interface {
	Close() error
	Set(key, value string) error
	SetWithTTL(key, value string, ttl time.Duration) error
	Read(key string) ([]byte, error)
}

type Badger struct {
	db *badger.DB
}

func CreateBadgerClient() (Client, error) {
	uuid := uuid2.New().String()
	badgerPath := filepath.Join("/tmp", "badger", uuid)
	badger, err := NewBadgerClient(badgerPath)
	if err != nil {
		return nil, err
	}

	return badger, nil
}

func NewInMemoryBadgerClient() (Client, error) {
	opt := badger.DefaultOptions("").WithInMemory(true)
	db, err := badger.Open(opt)
	if err != nil {
		return nil, err
	}

	badgerDb := Badger{db: db}

	return &badgerDb, nil
}

func NewBadgerClient(badgerPath string) (Client, error) {
	db, err := badger.Open(badger.DefaultOptions(badgerPath))
	if err != nil {
		return nil, err
	}

	badgerDb := Badger{db: db}

	return &badgerDb, nil
}
