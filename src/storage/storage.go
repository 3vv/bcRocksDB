package storage

import (
  "errors"
  . "../constants"
)

// File type
type FileType int

// File descriptor
type FileDesc struct {
	Type FileType
	Num  int64
}

// Common error
var (
	ErrInvalidFile = errors.New(PkgName+"/storage: invalid file for argument")
	ErrLocked      = errors.New(PkgName+"/storage: already locked")
	ErrClosed      = errors.New(PkgName+"/storage: closed")
)