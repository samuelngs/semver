package backend

import (
	"strings"

	"github.com/gocql/gocql"
	"github.com/samuelngs/semver/pkg/env"
)

var (
	cbegin = "BEGIN BATCH"
	cend   = "APPLY BATCH"
)

// Cassandra backend for semver
type Cassandra struct {
	*Core
	cluster *gocql.ClusterConfig
}

// Init method
func (c *Cassandra) Init() error {
	// cassandra databse hosts
	addr := strings.Split(
		env.Raw("SEMVER_BACKEND_ADDR", "localhost"),
		",", // split hosts incase environment is in 192.168.0.1,192.168.0.2 format
	)
	// database or keyspace name
	keyspace := env.Raw("SEMVER_BACKEND_DB", "semver")
	// create cassandra cluster client
	c.cluster = gocql.NewCluster(addr...)
	c.cluster.Keyspace = keyspace
	c.cluster.Consistency = gocql.Quorum
	c.cluster.ProtoVersion = 4
	session, err := c.cluster.CreateSession()
	if err != nil {
		return err
	}
	defer session.Close()
	// create keyspace and tables if they are not existed
	if err := session.Query(`CREATE KEYSPACE IF NOT EXISTS ? WITH replication = {'class': 'SimpleStrategy', 'replication_factor' : 3}`, keyspace).Exec(); err != nil {
		return err
	}
	if err := session.Query(`USE ?`, keyspace).Exec(); err != nil {
		return err
	}
	if err := session.Query(`CREATE TABLE IF NOT EXISTS db (id text, key text, val text, PRIMARY KEY (id, key))`).Exec(); err != nil {
		return err
	}
	return nil
}

// Name method
func (c *Cassandra) Name() string {
	return "cassandra"
}

// Path method
func (c *Cassandra) Path(key *Key) string {
	var dir string
	for i, str := range key.Dirs {
		if i > 0 {
			dir += "/"
		}
		dir += str
	}
	return dir
}

// Exists method
func (c *Cassandra) Exists(key *Key) (bool, error) {
	session, err := c.cluster.CreateSession()
	if err != nil {
		return false, err
	}
	defer session.Close()
	var count int
	if err := session.Query(`SELECT COUNT(*) FROM db WHERE id = ?`, key.ID).Consistency(gocql.One).Scan(&count); err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

// Set method
func (c *Cassandra) Set(v string, keys ...*Key) error {
	var query []string
	var args []interface{}
	var count int
	session, err := c.cluster.CreateSession()
	if err != nil {
		return err
	}
	defer session.Close()
	query = make([]string, len(keys)+2)
	query[count] = cbegin
	args = make([]interface{}, len(keys)*3)
	count++
	for i, key := range keys {
		k := c.Path(key)
		args[i*3] = key.ID
		args[i*3+1] = k
		args[i*3+2] = v
		query[count] = `INSERT INTO db (id, key, val) VALUES (?, ?, ?)`
		count++
	}
	query[count] = cend
	if err := session.Query(
		strings.Join(query, "\n"),
		args...,
	).Exec(); err != nil {
		return err
	}
	return nil
}

// Get method
func (c *Cassandra) Get(keys ...*Key) ([]string, error) {
	session, err := c.cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()
	res := []string{}
	ids := make(map[string][]string)
	for _, key := range keys {
		var keys []string
		if _, ok := ids[key.ID]; ok {
			keys = ids[key.ID]
		} else {
			keys = []string{}
		}
		k := c.Path(key)
		keys = append(keys, k)
		ids[key.ID] = keys
	}
	for i, o := range ids {
		var s string
		iter := session.Query(`SELECT val FROM db WHERE id = ? AND key in ?`, i, o).Iter()
		for iter.Scan(&s) {
			res = append(res, s)
		}
		if err := iter.Close(); err != nil {
			return nil, err
		}
	}
	return res, nil
}

// List method
func (c *Cassandra) List(key *Key) ([]*Key, error) {
	session, err := c.cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()
	keys := []*Key{}
	path := c.Path(key)
	iter := session.Query(`SELECT key FROM db WHERE id = ?`, key.ID).Iter()
	var k string
	for iter.Scan(&k) {
		if strings.HasPrefix(k, path) {
			keys = append(keys, &Key{ID: key.ID, Dirs: strings.Split(k, "/")})
		}
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return keys, nil
}

// Delete method
func (c *Cassandra) Delete(keys ...*Key) error {
	session, err := c.cluster.CreateSession()
	if err != nil {
		return err
	}
	defer session.Close()
	ids := make(map[string][]string)
	for _, key := range keys {
		if _, ok := ids[key.ID]; !ok {
			ids[key.ID] = []string{}
		}
		path := c.Path(key)
		ids[key.ID] = append(ids[key.ID], path)
	}
	queries := make([]string, len(ids)+2)
	args := make([]interface{}, len(ids)*2)
	count := 0
	queries[count] = cbegin
	queries[len(queries)-1] = cend
	for id, dirs := range ids {
		queries[count+1] = `DELETE FROM db WHERE id = ? AND key in ?`
		args[count*2] = id
		args[count*2+1] = dirs
		count++
	}
	if err := session.Query(
		strings.Join(queries, "\n"),
		args...,
	).Exec(); err != nil {
		return err
	}
	return nil
}
