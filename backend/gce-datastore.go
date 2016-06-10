package backend

import (
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
	}
	if e == nil {
		return false, nil
	}
	return true, nil
}

// Set method
func (d *GceDatastore) Set(val string, keys ...*Key) error {
	var e *Entity
	cache := map[string][]string{}
	for _, key := range keys {
		k := datastore.NewKey(d.ctx, "Semver", key.id, 0, nil)
		if err := d.client.Get(d.ctx, k, &e); err != nil && err != datastore.ErrNoSuchEntity {
			return err
		}
		if e.Archive == nil {
			e.Archive = make([]*Map, 0)
		}
		if cache[key.id] == nil {
			cache[key.id] = make([]string, len(e.Archive))
			for i, m := range e.Archive {
				cache[key.id][i] = m.Key
			}
		}
	}
	panic("you should override `set` method")
}

// Get method
func (d *GceDatastore) Get(keys ...*Key) ([]string, error) {
	panic("you should override `get` method")
}

// List method
func (d *GceDatastore) List(key *Key) ([]*Key, error) {
	panic("you should override `list` method")
}

// Delete method
func (d *GceDatastore) Delete(keys ...*Key) error {
	panic("you should override `delete` method")
}
