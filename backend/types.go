package backend

// Key struct
type Key struct {
	ID   string
	Dirs []string
}

// Entity represents a record
type Entity struct {
	Data string `json:"data"`
}

// Versioning represents a version record
type Versioning struct {
	Version string            `json:"version"`
	Archive map[string]string `json:"archive"`
}
