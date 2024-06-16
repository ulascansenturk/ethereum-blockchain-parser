package storage

import (
	"testing"
)

func TestInMemoryDatabase(t *testing.T) {
	db := New()

	t.Run("Test Put and Get", func(t *testing.T) {
		key := "testKey"
		value := [][]byte{[]byte("value1"), []byte("value2")}

		err := db.Put(key, value)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		got, err := db.Get(key)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(got) != len(value) || string(got[0]) != string(value[0]) || string(got[1]) != string(value[1]) {
			t.Fatalf("expected %v, got %v", value, got)
		}
	})

	t.Run("Test Has", func(t *testing.T) {
		key := "testKey"
		value := [][]byte{[]byte("value1"), []byte("value2")}

		err := db.Put(key, value)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		has, err := db.Has(key)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if !has {
			t.Fatalf("expected key to exist")
		}

		has, err = db.Has("nonexistentKey")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if has {
			t.Fatalf("expected key to not exist")
		}
	})

	t.Run("Test Get non-existent key", func(t *testing.T) {
		_, err := db.Get("nonexistentKey")
		if err != ErrNotFound {
			t.Fatalf("expected error %v, got %v", ErrNotFound, err)
		}
	})

	t.Run("Test operations on closed database", func(t *testing.T) {
		db.lock.Lock()
		db.db = nil
		db.lock.Unlock()

		err := db.Put("key", [][]byte{[]byte("value")})
		if err != ErrDBClosed {
			t.Fatalf("expected error %v, got %v", ErrDBClosed, err)
		}

		_, err = db.Get("key")
		if err != ErrDBClosed {
			t.Fatalf("expected error %v, got %v", ErrDBClosed, err)
		}

		_, err = db.Has("key")
		if err != ErrDBClosed {
			t.Fatalf("expected error %v, got %v", ErrDBClosed, err)
		}
	})
}
