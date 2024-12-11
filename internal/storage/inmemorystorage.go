package storage

import (
	"errors"
	"sync"
)

var (
	ErrNotFound = errors.New("key not found")
	ErrDBClosed = errors.New("database closed")
)

type InMemoryStorage interface {
	Has(key string) (bool, error)
	Get(key string) ([][]byte, error)
	Put(key string, value [][]byte) error
}

type Entry struct {
	Value [][]byte
}

type InMemoryDatabase struct {
	db   map[string]Entry
	lock sync.RWMutex
	done chan struct{}
	wg   sync.WaitGroup
}

func New() *InMemoryDatabase {
	db := &InMemoryDatabase{
		db:   make(map[string]Entry),
		done: make(chan struct{}),
	}
	db.wg.Add(1)
	return db
}

func (db *InMemoryDatabase) Has(key string) (bool, error) {
	db.lock.RLock()
	defer db.lock.RUnlock()

	if db.db == nil {
		return false, ErrDBClosed
	}
	_, ok := db.db[key]
	return ok, nil
}

func (db *InMemoryDatabase) Get(key string) ([][]byte, error) {
	db.lock.RLock()
	defer db.lock.RUnlock()

	if db.db == nil {
		return nil, ErrDBClosed
	}
	if entry, ok := db.db[key]; ok {
		return entry.Value, nil
	}
	return nil, ErrNotFound
}

func (db *InMemoryDatabase) Put(key string, value [][]byte) error {
	db.lock.Lock()
	defer db.lock.Unlock()

	if db.db == nil {
		return ErrDBClosed
	}

	db.db[key] = Entry{Value: value}
	return nil
}
