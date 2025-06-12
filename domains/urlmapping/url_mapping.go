package urlmapping

import "errors"

var (
	ErrNotFound = errors.New("The requested url was not found")
)

// TODO: Move interfaces here to prevent cyclic depdencies and promote pure depdency injection
