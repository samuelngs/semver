package backend

// New creates backend manager
func New(opts ...Client) *Manager {
	var c Client
	for _, o := range opts {
		c = o
		break
	}
	m := new(Manager)
	if c != nil {
		m.Use(c)
	}
	return m
}

// Manager represent the backend manager
type Manager struct {
	c Client
}

// Use to import backend client
func (m *Manager) Use(c Client) {
	m.c = c
	m.c.Init()
}

// before action
func (m *Manager) prepare() {
	if m.c == nil {
		panic("client has not been registered")
	}
}

// Name returns current client name
func (m *Manager) Name() string {
	m.prepare()
	return m.c.Name()
}

// Path creates a key
func (m *Manager) Path(id string, dirs ...string) *Key {
	m.prepare()
	return &Key{ID: id, Dirs: dirs}
}

// Exists checker
func (m *Manager) Exists(key *Key) (bool, error) {
	m.prepare()
	return m.c.Exists(key)
}

// Set data to storage
func (m *Manager) Set(val string, keys ...*Key) error {
	m.prepare()
	return m.c.Set(val, keys...)
}

// Get method
func (m *Manager) Get(keys ...*Key) ([]string, error) {
	m.prepare()
	return m.c.Get(keys...)
}

// List method
func (m *Manager) List(key *Key) ([]*Key, error) {
	m.prepare()
	return m.c.List(key)
}

// Delete method
func (m *Manager) Delete(keys ...*Key) error {
	m.prepare()
	return m.c.Delete(keys...)
}
