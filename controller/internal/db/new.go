package db

import (
	"log"

	"github.com/boltdb/bolt"
)

type DB struct {
	bolt *bolt.DB
}

func New(addr string) *DB {
	db, err := bolt.Open(addr, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("main"))
		return err
	})

	if err != nil {
		log.Fatal(err)
	}
	return &DB{bolt: db}
}

func (d *DB) Close() error {
	return d.bolt.Close()
}

func (d *DB) Set(key string, value []byte) error {
	return d.bolt.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("main"))
		return bucket.Put([]byte(key), value)
	})
}

func (d *DB) Get(key string) ([]byte, error) {
	var value []byte
	err := d.bolt.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("main"))
		value = bucket.Get([]byte(key))
		return nil
	})
	return value, err
}
