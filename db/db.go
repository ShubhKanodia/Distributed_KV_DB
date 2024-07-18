package db

import (
	"fmt"

	bolt "go.etcd.io/bbolt"
)

//first pipe in our package

// Database - open bolt database

var defaultBucket = []byte("default")

type Database struct { //key-value store
	db *bolt.DB
}

func NewDatabase(dbPath string) (db *Database, closeFunc func() error, err error) {
	boltDb, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		return nil, nil, err
	}

	//not required we use close
	// cleanupFunc : = boltDb.close
	// var cleanupFunc func() error
	// defer func() {
	// 	if cleanupFunc != nil {
	// 		cleanupFunc()

	// 	}
	// }()

	db = &Database{db: boltDb}
	closeFunc = boltDb.Close

	// cleanupFunc = closeFunc
	if err := db.createDefaultBucket(); err != nil {
		closeFunc()
		return nil, nil, fmt.Errorf("creating default bucket: %w", err)
	}

	return db, closeFunc, nil
}

func (d *Database) createDefaultBucket() error {
	return d.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(defaultBucket)
		return err
	})
}

// SetKey sets the key to the requested value into the default database or returns an error.
func (d *Database) SetKey(key string, value []byte) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(defaultBucket)
		return b.Put([]byte(key), value)
	})
}

// GetKey get the value of the requested from a default database.
func (d *Database) GetKey(key string) ([]byte, error) {
	var result []byte
	err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(defaultBucket)
		result = b.Get([]byte(key))
		return nil
	})

	if err == nil {
		return result, nil
	}
	return nil, err
}
