package db

import (
	"encoding/binary"
	"time"

	"github.com/boltdb/bolt"
)

const dbPath = "/Users/ruslanmuradov/github.com:gasuhwbab/cli_todo_app/data/todo_app.db"

type Storage struct {
	db   *bolt.DB
	path string
}

func NewStorage(path string) *Storage {
	return &Storage{path: path}
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

func (storage *Storage) Add(buf []byte) error {
	if err := storage.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(time.Now().Format("2006/01/02")))
		if err != nil {
			return err
		}
		if err := bucket.Put(buf[:2], buf[2:]); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (storage *Storage) Update(from, to []byte) error {
	if err := storage.db.Update(func(tx *bolt.Tx) error {
		id := from[0:2]
		time := time.Unix(0, int64(binary.BigEndian.Uint64(from[2:10]))).Format("2006/01/02")
		bucket := tx.Bucket([]byte(time))
		if err := bucket.Put(id, to[2:]); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (storage *Storage) Delete(buf []byte) error {
	if err := storage.db.Update(func(tx *bolt.Tx) error {
		id := buf[0:2]
		time := time.Unix(0, int64(binary.BigEndian.Uint64(buf[2:10]))).Format("2006/01/02")
		bucket := tx.Bucket([]byte(time))
		if err := bucket.Delete(id); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (storage *Storage) Get() ([][]byte, error) {
	bufs := make([][]byte, 0)
	storage.db.View(func(tx *bolt.Tx) error {
		tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			b.ForEach(func(k, v []byte) error {
				bufs = append(bufs, v)
				return nil
			})
			return nil
		})
		return nil
	})
	return bufs, nil
}
