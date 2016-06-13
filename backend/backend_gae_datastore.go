package backend

import (
	"encoding/json"
	"strings"

	"github.com/samuelngs/semver/pkg/env"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"

	"google.golang.org/appengine"
	"google.golang.org/cloud"
	"google.golang.org/cloud/datastore"
)

// GceDatastore backend for semver
type GceDatastore struct {
	*Core
}

// client creates gce client
func (d *GceDatastore) storage() (context.Context, *datastore.Client, error) {
	key := []byte(env.Raw("SEMVER_BACKEND_TOKEN"))
	conf, err := google.JWTConfigFromJSON(key, datastore.ScopeDatastore)
	if err != nil {
		return nil, nil, err
	}
	ctx := appengine.BackgroundContext()
	client, err := datastore.NewClient(ctx, appengine.AppID(ctx), cloud.WithTokenSource(conf.TokenSource(ctx)))
	if err != nil {
		return nil, nil, err
	}
	return ctx, client, nil
}

// Init method
func (d *GceDatastore) Init() error {
	return nil
}

// Name method
func (d *GceDatastore) Name() string {
	return "gce-datastore"
}

// Path method
func (d *GceDatastore) Path(key *Key) string {
	var dir string
	for i, str := range key.Dirs {
		if i > 0 {
			dir += ":"
		}
		dir += str
	}
	return dir
}

// Exists method
func (d *GceDatastore) Exists(key *Key) (bool, error) {
	var e Entity
	ctx, client, err := d.storage()
	if err != nil {
		return false, err
	}
	k := datastore.NewKey(ctx, "Semver", key.ID, 0, nil)
	if err := client.Get(ctx, k, &e); err != nil && err != datastore.ErrNoSuchEntity {
		return false, err
	} else if err != nil && err == datastore.ErrNoSuchEntity {
		return false, nil
	}
	return true, nil
}

// Set method
func (d *GceDatastore) Set(val string, keys ...*Key) error {
	ctx, client, err := d.storage()
	if err != nil {
		return err
	}
	// entity cache
	cache := make(map[string]*Entity)
	for _, key := range keys {
		var e *Entity
		var v *Versioning
		if o, ok := cache[key.ID]; ok {
			e = o
		} else {
			k := datastore.NewKey(ctx, "Semver", key.ID, 0, nil)
			if err := client.Get(ctx, k, &e); err != nil && err != datastore.ErrNoSuchEntity {
				return err
			}
		}
		if e == nil {
			e = new(Entity)
		}
		if e.Data != "" {
			if err := json.Unmarshal([]byte(e.Data), &v); err != nil {
				return err
			}
		} else {
			v = &Versioning{Archive: make(map[string]string)}
		}
		p := d.Path(key)
		if p == "version" {
			v.Version = val
		} else {
			v.Archive[p] = val
		}
		b, err := json.Marshal(v)
		if err != nil {
			return err
		}
		e.Data = string(b[:])
		cache[key.ID] = e
	}
	for i, e := range cache {
		k := datastore.NewKey(ctx, "Semver", i, 0, nil)
		if _, err := client.Put(ctx, k, e); err != nil {
			return err
		}
	}
	return nil
}

// Get method
func (d *GceDatastore) Get(keys ...*Key) ([]string, error) {
	ctx, client, err := d.storage()
	if err != nil {
		return nil, err
	}
	// entity cache
	cache := make(map[string]*Entity)
	// return values
	store := []string{}
	for _, key := range keys {
		var e *Entity
		var v *Versioning
		if o, ok := cache[key.ID]; ok {
			e = o
		} else {
			k := datastore.NewKey(ctx, "Semver", key.ID, 0, nil)
			if err := client.Get(ctx, k, &e); err != nil && err != datastore.ErrNoSuchEntity {
				return nil, err
			} else if err != nil && err == datastore.ErrNoSuchEntity {
				return nil, ErrRecordNotFound
			}
			if err := json.Unmarshal([]byte(e.Data), &v); err != nil {
				return nil, err
			}
			cache[key.ID] = e
		}
		p := d.Path(key)
		if p == "version" {
			store = append(store, v.Version)
		} else {
			for k, s := range v.Archive {
				if p == k {
					store = append(store, s)
					break
				}
			}
		}
	}
	return store, nil
}

// List method
func (d *GceDatastore) List(key *Key) ([]*Key, error) {
	var e *Entity
	var v *Versioning
	ctx, client, err := d.storage()
	if err != nil {
		return nil, err
	}
	r := []*Key{}
	k := datastore.NewKey(ctx, "Semver", key.ID, 0, nil)
	if err := client.Get(ctx, k, &e); err != nil && err != datastore.ErrNoSuchEntity {
		return nil, err
	} else if err != nil && err == datastore.ErrNoSuchEntity {
		return nil, ErrRecordNotFound
	}
	if err := json.Unmarshal([]byte(e.Data), &v); err != nil {
		return nil, err
	}
	for k := range v.Archive {
		s := strings.Split(k, ":")
		r = append(r, &Key{ID: key.ID, Dirs: s})
	}
	return r, nil
}

// Delete method
func (d *GceDatastore) Delete(keys ...*Key) error {
	ctx, client, err := d.storage()
	if err != nil {
		return err
	}
	// entity cache
	cache := make(map[string]*Entity)
	for _, key := range keys {
		var e *Entity
		var v *Versioning
		if o, ok := cache[key.ID]; ok {
			e = o
		} else {
			k := datastore.NewKey(ctx, "Semver", key.ID, 0, nil)
			if err := client.Get(ctx, k, &v); err != nil && err != datastore.ErrNoSuchEntity {
				return err
			} else if err != nil && err == datastore.ErrNoSuchEntity {
				continue
			}
			if err := json.Unmarshal([]byte(e.Data), &v); err != nil {
				return err
			}
		}
		p := d.Path(key)
		if p == "version" {
			v.Version = ""
		} else {
			delete(v.Archive, p)
		}
		b, err := json.Marshal(v)
		if err != nil {
			return err
		}
		e.Data = string(b[:])
		cache[key.ID] = e
	}
	for i, e := range cache {
		var v *Versioning
		k := datastore.NewKey(ctx, "Semver", i, 0, nil)
		if err := json.Unmarshal([]byte(e.Data), &v); err != nil {
			return err
		}
		if len(v.Archive) > 0 {
			if _, err := client.Put(ctx, k, e); err != nil {
				return err
			}
		} else {
			if err := client.Delete(ctx, k); err != nil {
				return err
			}
		}
	}
	return nil
}
