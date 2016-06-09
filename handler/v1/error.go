package v1

import "errors"

// List of error messages
var (
	ErrProjectNotFound         = errors.New("project id does not match any records in our database")
	ErrInvalidVersioningFormat = errors.New("invalid semantic versioning format")
	ErrInvalidUUID             = errors.New("invalid uuid")
)
