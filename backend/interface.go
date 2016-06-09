package backend

// Client interface
type Client interface {
	Init() error
	Name() string
	Path(key *Key) string
	Exists(key *Key) (bool, error)
	Set(val string, keys ...*Key) error
	Get(keys ...*Key) ([]string, error)
	List(key *Key) ([]*Key, error)
	Delete(keys ...*Key) error
}

// Core for extend purpose
type Core struct{}

// Init method
func (e *Core) Init() string {
	panic("you should override `init` method")
}

// Name method
func (e *Core) Name() string {
	panic("you should override `name` method")
}

// Path method
func (e *Core) Path(key *Key) string {
	panic("you should override `path` method")
}

// Exists method
func (e *Core) Exists(key *Key) (bool, error) {
	panic("you should override `exists` method")
}

// Set method
func (e *Core) Set(val string, keys ...*Key) error {
	panic("you should override `set` method")
}

// Get method
func (e *Core) Get(keys ...*Key) ([]string, error) {
	panic("you should override `get` method")
}

// List method
func (e *Core) List(key *Key) ([]*Key, error) {
	panic("you should override `list` method")
}

// Delete method
func (e *Core) Delete(keys ...*Key) error {
	panic("you should override `delete` method")
}
