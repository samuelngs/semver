// +build !appengine

package backend

import (
	"bytes"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/samuelngs/semver/pkg/env"
)

// Bolt backend for semver
type Bolt struct {
	*Core
	db *bolt.DB
}

// Init method
func (b *Bolt) Init() error {
	db, err := bolt.Open(
		env.Raw("SEMVER_BACKEND_ADDR", "local.db"),
		0600,
		nil,
	)
	if err != nil {
		return err
	}
	b.db = db
	return nil
}

// Name method
func (b *Bolt) Name() string {
	return "bolt"
}

// Path method
func (b *Bolt) Path(key *Key) string {
	var dir string
	for i, str := range key.dirs {
		if i > 0 {
			dir += ":"
		}
		dir += str
	}
	return dir
}

// Exists method
func (b *Bolt) Exists(key *Key) (bool, error) {
	var exists bool
	err := b.db.View(func(tx *bolt.Tx) error {
		if bucket := tx.Bucket([]byte(key.id)); bucket != nil {
			exists = true
		}
		return nil
	})
	if err != nil {
		return false, err
	}
	return exists, nil
}

// Set method
func (b *Bolt) Set(val string, keys ...*Key) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		for _, key := range keys {
			bucket, err := tx.CreateBucketIfNotExists([]byte(key.id))
			if err != nil {
				return err
			}
			if err := bucket.Put(
				[]byte(b.Path(key)),
				[]byte(val),
			); err != nil {
				return err
			}
		}
		return nil
	})
}

// Get method
func (b *Bolt) Get(keys ...*Key) ([]string, error) {
	vals := []string{}
	err := b.db.View(func(tx *bolt.Tx) error {
		for _, key := range keys {
			bucket := tx.Bucket([]byte(key.id))
			if bucket == nil {
				return ErrRecordNotFound
			}
			v := bucket.Get([]byte(b.Path(key)))
			vals = append(vals, string(v[:]))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return vals, nil
}

// List method
func (b *Bolt) List(key *Key) ([]*Key, error) {
	keys := []*Key{}
	err := b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(key.id))
		if bucket == nil {
			return ErrRecordNotFound
		}
		cursor := bucket.Cursor()
		prefix := []byte(b.Path(key))
		if len(prefix) > 0 {
			for k, _ := cursor.Seek(prefix); bytes.HasPrefix(k, prefix); k, _ = cursor.Next() {
				keys = append(keys, &Key{id: key.id, dirs: strings.Split(string(k[:]), ":")})
			}
		} else {
			for k, _ := cursor.First(); k != nil; k, _ = cursor.Next() {
				keys = append(keys, &Key{id: key.id, dirs: strings.Split(string(k[:]), ":")})
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return keys, err
}

// Delete method
func (b *Bolt) Delete(keys ...*Key) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		for _, key := range keys {
			if path := b.Path(key); path == "" {
				if err := tx.DeleteBucket([]byte(key.id)); err != nil {
					return err
				}
			} else {
				var count int
				bucket := tx.Bucket([]byte(key.id))
				if bucket == nil {
					continue
				}
				if err := bucket.Delete([]byte(path)); err != nil {
					return err
				}
				bucket.ForEach(func(_, _ []byte) error {
					count++
					return nil
				})
				if count == 0 {
					if err := tx.DeleteBucket([]byte(key.id)); err != nil {
						return err
					}
				}
			}
		}
		return nil
	})
}
