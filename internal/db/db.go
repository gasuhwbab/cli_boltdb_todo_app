package db

import (
	"encoding/binary"
	"errors"

	"github.com/boltdb/bolt"
)

const dbPath = "/Users/ruslanmuradov/github.com:gasuhwbab/cli_todo_app/data/todo_app.db"

var Db = &Storage{path: dbPath}

type Storage struct {
	db   *bolt.DB
	path string
}

func (storage *Storage) StartDb() error {
	db, err := bolt.Open(storage.path, 0600, nil)
	if err != nil {
		return err
	}
	storage.db = db
	return nil
}

func (storage *Storage) Close() error {
	return storage.db.Close()
}

func (storage *Storage) Add(name, buf []byte) error {
	if err := storage.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(name)

		if err != nil {
			return err
		}
		id, _ := bucket.NextSequence()
		binaryId := make([]byte, 8)
		binary.BigEndian.PutUint64(binaryId, id)
		if err := bucket.Put(binaryId, buf[8:]); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (storage *Storage) Update(name, id, to []byte) error {
	if err := storage.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(name)
		if err := bucket.Put(id, to[8:]); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (storage *Storage) Delete(name, id []byte) error {
	if err := storage.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(name)
		if err := bucket.Delete(id); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (storage *Storage) Get(name []byte) ([][]byte, error) {
	bufs := make([][]byte, 0)
	if err := storage.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(name)
		if bucket == nil {
			return errors.New("bucket does not exist")
		}
		c := bucket.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			buf := make([]byte, len(k)+len(v))
			copy(buf[:len(k)], k)
			copy(buf[len(k):], v)
			bufs = append(bufs, buf)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return bufs, nil
}
