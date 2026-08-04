package archytas

import (
	"path/filepath"
	"time"

	"github.com/dgraph-io/badger/v4"
	uuid2 "github.com/google/uuid"
)

type Client interface {
	Close() error
	Set(key, value string) error
	SetBytes(key string, value []byte) error
	SetWithTTL(key, value string, ttl time.Duration) error
	SetBytesWithTTL(key string, value []byte, ttl time.Duration) error
	Read(key string) ([]byte, error)
	Get(key string) ([]byte, error)
	GetString(key string) (string, error)
	Exists(key string) (bool, error)
	Delete(key string) error
	DeletePrefix(prefix string) error
	Stats() (Stats, error)
}

type Badger struct {
	db *badger.DB
}

// Option changes the Badger options used when opening a client.
type Option func(badger.Options) badger.Options

// WithLogging controls Badger's built-in logger. Logging is enabled by default.
func WithLogging(enabled bool) Option {
	return func(options badger.Options) badger.Options {
		if !enabled {
			return options.WithLogger(nil)
		}
		return options
	}
}

// WithLogger installs a custom Badger logger. A nil logger disables logging.
func WithLogger(logger badger.Logger) Option {
	return func(options badger.Options) badger.Options {
		return options.WithLogger(logger)
	}
}

func CreateBadgerClient() (Client, error) {
	return CreateBadgerClientWithOptions()
}

// CreateBadgerClientWithOptions creates a client in a unique temporary
// directory. The database remains on disk after Close.
func CreateBadgerClientWithOptions(options ...Option) (Client, error) {
	uuid := uuid2.New().String()
	badgerPath := filepath.Join("/tmp", "badger", uuid)
	badger, err := NewBadgerClientWithOptions(badgerPath, options...)
	if err != nil {
		return nil, err
	}

	return badger, nil
}

func NewInMemoryBadgerClient() (Client, error) {
	return NewInMemoryBadgerClientWithOptions()
}

// NewInMemoryBadgerClientWithOptions opens an in-memory client with optional
// configuration such as WithLogging(false).
func NewInMemoryBadgerClientWithOptions(options ...Option) (Client, error) {
	return openBadger(applyOptions(badger.DefaultOptions("").WithInMemory(true), options))
}

func NewBadgerClient(badgerPath string) (Client, error) {
	return NewBadgerClientWithOptions(badgerPath)
}

// NewBadgerClientWithOptions opens a disk-backed client with optional
// configuration such as WithLogging(false) or WithLogger(logger).
func NewBadgerClientWithOptions(badgerPath string, options ...Option) (Client, error) {
	return openBadger(applyOptions(badger.DefaultOptions(badgerPath), options))
}

func applyOptions(options badger.Options, modifiers []Option) badger.Options {
	for _, modify := range modifiers {
		if modify != nil {
			options = modify(options)
		}
	}
	return options
}

func openBadger(options badger.Options) (Client, error) {
	db, err := badger.Open(options)
	if err != nil {
		return nil, err
	}

	return &Badger{db: db}, nil
}
