package kvdb

import (
	"schedrestd/common"
	"schedrestd/config"
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/boltdb/bolt"
	"go.uber.org/fx"
)

// Module Module
var Module = fx.Options(
	fx.Provide(NewKVStore),
)

// KVStore provides k/v Insert/Get/Update/Delete interface
type KVStore struct {
	dbPath string
	db     *bolt.DB
}

// KeyError will be thrown when queried key not found
type KeyError struct {
	k string
}

// TableError will be thrown when query table not found
type TableError struct {
	t string
}

// Error Error
func (ke *KeyError) Error() string {
	return fmt.Sprintf("key %s not found", ke.k)
}

// Error Error
func (te *TableError) Error() string {
	return fmt.Sprintf("table %s not found", te.t)
}

// str2Byte is to transform string to bytes
func str2Byte(str string) []byte {
	return []byte(str)
}

// NewKVStore creates a kv db
func NewKVStore(conf *config.Config) *KVStore {
	return &KVStore{
		dbPath: fmt.Sprintf("%v/%v", conf.WorkDir, common.BoltDBName),
	}
}

// Open opens the store.
// KV DB must be opened before it is used.
func (s *KVStore) Open() error {
	db, err := bolt.Open(s.dbPath, 0600, nil)
	if err != nil {
		return err
	}
	s.db = db
	return nil
}

// Close closes the db
func (s *KVStore) Close() {
	s.db.Close()
}

// CreateTableIfNotExist creates a table if it doesn't exist
// nothing will happen if the table has existed
// Thread-Safe
func (s *KVStore) CreateTableIfNotExist(table string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(str2Byte(table))
		return err
	})
}

// DeleteTable deletes a table
// Returns an error if the table cannot be found or if the key represents a non-table value.
// Thread-Safe
func (s *KVStore) DeleteTable(table string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket(str2Byte(table))
	})
}

// HasTable checks whether a table exists
func (s *KVStore) HasTable(table string) bool {
	var has bool
	s.db.View(func(tx *bolt.Tx) error {
		has = tx.Bucket(str2Byte(table)) != nil
		return nil
	})
	return has
}

// Put Put
// Insert k/v in table
// Thread Safe
func (s *KVStore) Put(table, k string, v []byte) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(str2Byte(table))
		if b == nil {
			return &TableError{table}
		}
		return b.Put(str2Byte(k), v)
	})
}

// Get gets the value of k in table
// Thread-Safe
func (s *KVStore) Get(table, k string) ([]byte, error) {
	var value []byte
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(str2Byte(table))
		if b == nil {
			return &TableError{table}
		}

		v := b.Get(str2Byte(k))
		if v == nil {
			return &KeyError{k}
		}
		value = make([]byte, len(v))
		copy(value, v)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return value, nil
}

// Find query the values with prefix k in table
// Thread-Safe
func (s *KVStore) Find(table, k string) ([]string, error) {
	var values []string
	err := s.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(str2Byte(table)).Cursor()
		if c == nil {
			return &TableError{table}
		}

		prefix := []byte(k)
		for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
			values = append(values, string(v))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return values, nil
}

// DeleteKey deletes key from table
// Thread-Safe
func (s *KVStore) DeleteKey(table, key string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(str2Byte(table))
		if bucket == nil {
			return &TableError{table}
		}
		return bucket.Delete(str2Byte(key))
	})
}

// RemoveDB will clear the database source file
func (s *KVStore) RemoveDB() error {
	// store the dbPath because path in s.db will be reset after close
	dbPath := s.db.Path()
	s.Close()
	return os.Remove(dbPath)
}

// RemoveSmallValues remove the smaller key/value store comparing with specified value
func (s *KVStore) RemoveSmallValues(table, value string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(str2Byte(table))
		if bucket == nil {
			return &TableError{table}
		}
		var keys []string
		bucket.ForEach(func(k, v []byte) error {
			if strings.Compare(string(v), value) == -1 {
				keys = append(keys, string(k))
			}
			return nil
		})

		for _, key := range keys {
			err := bucket.Delete(str2Byte(key))
			if err != nil {
				return err
			}
		}

		return nil
	})
}
