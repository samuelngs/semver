package backend

import (
	"fmt"
	"strings"

	"gopkg.in/redis.v3"

	"github.com/samuelngs/semver/pkg/env"
)

// Redis backend for semver
type Redis struct {
	*Core
	c *redis.Client
}

// Init method
func (r *Redis) Init() error {
	opts := &redis.Options{
		Addr:       env.Raw("SEMVER_BACKEND_ADDR", "localhost:6379"),
		DB:         env.I64("SEMVER_BACKEND_DB", 0),
		MaxRetries: env.Int("SEMVER_BACKEND_RETRIES", 5),
	}
	r.c = redis.NewClient(opts)
	return nil
}

// Name method
func (r *Redis) Name() string {
	return "redis"
}

// Path method
func (r *Redis) Path(key *Key) string {
	dir := fmt.Sprintf("semver:db:%s", key.ID)
	for _, str := range key.Dirs {
		dir += ":" + str
	}
	return dir
}

// Exists method
func (r *Redis) Exists(key *Key) (bool, error) {
	_, err := r.c.Get(r.Path(key)).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

// Set method
func (r *Redis) Set(val string, keys ...*Key) error {
	items := make([]string, len(keys)*2)
	for i, key := range keys {
		id := r.Path(key)
		items[i*2] = id
		items[i*2+1] = val
	}
	if err := r.c.MSet(items...).Err(); err != nil {
		return err
	}
	return nil
}

// Get method
func (r *Redis) Get(keys ...*Key) ([]string, error) {
	dirs := make([]string, len(keys))
	for i, key := range keys {
		dirs[i] = r.Path(key)
	}
	vs, err := r.c.MGet(dirs...).Result()
	res := make([]string, len(vs))
	for i, s := range vs {
		res[i] = s.(string)
	}
	if err != nil {
		return nil, err
	}
	return res, nil
}

// List method
func (r *Redis) List(key *Key) ([]*Key, error) {
	temp := &Key{ID: key.ID, Dirs: make([]string, len(key.Dirs))}
	for i, v := range key.Dirs {
		temp.Dirs[i] = v
	}
	if len(temp.Dirs) > 0 {
		if s := temp.Dirs[len(temp.Dirs)-1]; !strings.Contains(s, "*") {
			temp.Dirs = append(temp.Dirs, "*")
		}
	} else {
		temp.Dirs = append(temp.Dirs, "*")
	}
	tar := strings.TrimSuffix(r.Path(temp), ":*")
	ids, err := r.c.Keys(r.Path(temp)).Result()
	if err != nil {
		return nil, err
	}
	keys := make([]*Key, len(ids))
	for i, id := range ids {
		pfx := strings.Replace(id, tar, "", -1)
		str := strings.TrimPrefix(pfx, ":")
		parts := strings.SplitAfter(str, ":")
		dirs := make([]string, len(key.Dirs)+len(parts))
		for i, v := range key.Dirs {
			dirs[i] = v
		}
		for i, v := range parts {
			dirs[i+len(key.Dirs)] = v
		}
		keys[i] = &Key{
			ID:   key.ID,
			Dirs: dirs,
		}
	}
	return keys, nil
}

// Delete method
func (r *Redis) Delete(keys ...*Key) error {
	ids := make([]string, len(keys))
	for i, v := range keys {
		ids[i] = r.Path(v)
	}
	if _, err := r.c.Del(ids...).Result(); err != nil {
		return err
	}
	return nil
}
