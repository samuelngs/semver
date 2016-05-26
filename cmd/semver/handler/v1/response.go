package v1

// Err represent error message
type Err struct {
	Error string `json:"error" xml:"message"`
}

// Response represents a valid semver version
type Response struct {
	ID      string   `json:"id,omitempty" xml:"id,omitempty"`
	Version string   `json:"version,omitempty" xml:"version,omitempty"`
	Major   uint64   `json:"major" xml:"major"`
	Minor   uint64   `json:"minor" xml:"minor"`
	Patch   uint64   `json:"patch" xml:"patch"`
	Build   []string `json:"build,omitempty" xml:"build,omitempty"`
}

// List represents a list of semver version
type List struct {
	Versions []interface{} `json:"versions" xml:"version"`
}
