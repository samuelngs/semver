package backend

import "errors"

// List of error messages
var (
	ErrRecordNotFound = errors.New("does not match any records in our database")
)
