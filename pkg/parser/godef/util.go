package godef

import "errors"

var (
	ErrNotExported = errors.New("type is not exported")
)
