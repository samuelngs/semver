package backend

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
