package backend

// list of available backends
var backends = []Client{
	new(Bolt),
	new(Redis),
	new(GceDatastore),
}

// Get by name
func Get(s string) Client {
	for _, i := range backends {
		if i.Name() == s {
			return i
		}
	}
	return backends[0]
}

// Key struct
type Key struct {
	id   string
	dirs []string
}

// Entity represents a version record
type Entity struct {
	Version string
	Archive []*Map
}

// Map represents map type object
type Map struct {
	Key string
	Val string
}
