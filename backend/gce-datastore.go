package backend

import (
	"strings"

	"github.com/samuelngs/semver/pkg/env"

	"golang.org/x/net/context"
	"google.golang.org/cloud/datastore"
)

// GceDatastore backend for semver
type GceDatastore struct {
	*Core
	ctx    context.Context
	client *datastore.Client
}

// Init method
func (d *GceDatastore) Init() error {
	d.ctx = context.Background()
	client, err := datastore.NewClient(
		d.ctx,
		env.Raw("SEMVER_BACKEND_ADDR", "semver-co"),
	)
	if err != nil {
		return err
	}
	d.client = client
	return nil
}

// Name method
func (d *GceDatastore) Name() string {
	return "gce-datastore"
}

// Path method
func (d *GceDatastore) Path(key *Key) string {
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
func (d *GceDatastore) Exists(key *Key) (bool, error) {
	var e *Entity
	k := datastore.NewKey(d.ctx, "Semver", key.id, 0, nil)
	if err := d.client.Get(d.ctx, k, &e); err != nil && err != datastore.ErrNoSuchEntity {
		return false, err
	} else if err != nil && err == datastore.ErrNoSuchEntity {
		return false, nil
	}
	return true, nil
}

// Set method
func (d *GceDatastore) Set(val string, keys ...*Key) error {
	// entity cache
	cache := make(map[string]*Entity)
	for _, key := range keys {
		var e *Entity
		if o, ok := cache[key.id]; ok {
			e = o
		} else {
			k := datastore.NewKey(d.ctx, "Semver", key.id, 0, nil)
			if err := d.client.Get(d.ctx, k, &e); err != nil && err != datastore.ErrNoSuchEntity {
				return err
			}
		}
		p := d.Path(key)
		if p == "version" {
			e.Version = val
		} else {
			x := -1
			for i, m := range e.Archive {
				if p == m.Key {
					x = i
					break
				}
			}
			if x > -1 {
				e.Archive[x].Val = val
			} else {
				e.Archive = append(e.Archive, &Map{Key: p, Val: val})
			}
		}
		cache[key.id] = e
	}
	for i, e := range cache {
		k := datastore.NewKey(d.ctx, "Semver", i, 0, nil)
		if _, err := d.client.Put(d.ctx, k, e); err != nil {
			return err
		}
	}
	return nil
}

// Get method
func (d *GceDatastore) Get(keys ...*Key) ([]string, error) {
	// entity cache
	cache := make(map[string]*Entity)
	// return values
	store := []string{}
	for _, key := range keys {
		var e *Entity
		if o, ok := cache[key.id]; ok {
			e = o
		} else {
			k := datastore.NewKey(d.ctx, "Semver", key.id, 0, nil)
			if err := d.client.Get(d.ctx, k, &e); err != nil && err != datastore.ErrNoSuchEntity {
				return nil, err
			} else if err != nil && err == datastore.ErrNoSuchEntity {
				return nil, ErrRecordNotFound
			}
			cache[key.id] = e
		}
		p := d.Path(key)
		if p == "version" {
			store = append(store, e.Version)
		} else {
			var s string
			for _, m := range e.Archive {
				if p == m.Key {
					s = m.Val
					break
				}
			}
			store = append(store, s)
		}
	}
	return store, nil
}

// List method
func (d *GceDatastore) List(key *Key) ([]*Key, error) {
	var e *Entity
	r := []*Key{}
	k := datastore.NewKey(d.ctx, "Semver", key.id, 0, nil)
	if err := d.client.Get(d.ctx, k, &e); err != nil && err != datastore.ErrNoSuchEntity {
		return nil, err
	} else if err != nil && err == datastore.ErrNoSuchEntity {
		return nil, ErrRecordNotFound
	}
	for _, m := range e.Archive {
		s := strings.Split(m.Key, ":")
		r = append(r, &Key{id: key.id, dirs: s})
	}
	return r, nil
}

// Delete method
func (d *GceDatastore) Delete(keys ...*Key) error {
	// entity cache
	cache := make(map[string]*Entity)
	for _, key := range keys {
		var e *Entity
		if o, ok := cache[key.id]; ok {
			e = o
		} else {
			k := datastore.NewKey(d.ctx, "Semver", key.id, 0, nil)
			if err := d.client.Get(d.ctx, k, &e); err != nil && err != datastore.ErrNoSuchEntity {
				return err
			} else if err != nil && err == datastore.ErrNoSuchEntity {
				continue
			}
		}
		p := d.Path(key)
		if p == "version" {
			e.Version = ""
		} else {
			x := -1
			for i, m := range e.Archive {
				if p == m.Key {
					x = i
					break
				}
			}
			if x > -1 {
				e.Archive = append(e.Archive[:x], e.Archive[x+1:]...)
			}
		}
		cache[key.id] = e
	}
	for i, e := range cache {
		k := datastore.NewKey(d.ctx, "Semver", i, 0, nil)
		if len(e.Archive) > 0 {
			if _, err := d.client.Put(d.ctx, k, e); err != nil {
				return err
			}
		} else {
			if err := d.client.Delete(d.ctx, k); err != nil {
				return err
			}
		}
	}
	return nil
}
