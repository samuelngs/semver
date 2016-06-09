package v1

import "fmt"

// Warning represent error message
type Warning struct {
	Error string `json:"error" xml:"message"`
}

// Versioning represents a valid semver version
type Versioning struct {
	Project string   `json:"project,omitempty" xml:"project,omitempty"`
	Version string   `json:"version,omitempty" xml:"version,omitempty"`
	Major   uint64   `json:"major" xml:"major"`
	Minor   uint64   `json:"minor" xml:"minor"`
	Patch   uint64   `json:"patch" xml:"patch"`
	Build   []string `json:"build,omitempty" xml:"build,omitempty"`
}

// String returns the string format of Versioning object
func (v *Versioning) String() string {
	if v.Project != "" {
		return v.Project
	}
	return v.Version
}

// Archive represents a list of semver version
type Archive struct {
	Versions []*Versioning `json:"versions" xml:"version"`
}

// String returns the string format of Archive object
func (v *Archive) String() string {
	var output string
	if v.Versions == nil {
		v.Versions = make([]*Versioning, 0)
	}
	for _, ver := range v.Versions {
		output += fmt.Sprintf("%v\n", ver)
	}
	return output
}
