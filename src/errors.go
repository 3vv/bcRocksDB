package rocksdb

import (
	"errors"
	. "./constants"
)

// Common errors (in alphabetical order)
var (
	ErrNotFound         = errors.New(PkgName + ": not found")
)
