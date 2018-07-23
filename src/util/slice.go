package util

//#include <stdlib.h>
import "C"
import (
	"unsafe"
)

// Slice is used as a wrapper for non-copy values
type Slice struct {
	Data_  *C.char
	Size_  C.size_t
	Freed_ bool
}

// NewSlice returns a slice with the given data.
func NewSlice(data *C.char, size C.size_t) *Slice {
	return &Slice{data, size, false}
}

func NewSliceFreed(data *C.char, size C.size_t, freed bool) *Slice {
	return &Slice{data, size, freed}
}

// StringToSlice is similar to NewSlice, but can be called with
// a Go string type. This exists to make testing integration
// with Gorocksdb easier.
func StringToSlice(data string) *Slice {
	return NewSlice(C.CString(data), C.size_t(len(data)))
}

// Data returns the data of the slice.
func (s *Slice) Data() []byte {
	return charToByte(s.Data_, s.Size_)
}

// Size returns the size of the data.
func (s *Slice) Size() int {
	return int(s.Size_)
}

// Free frees the slice data.
func (s *Slice) Free() {
	if !s.Freed_ {
		C.free(unsafe.Pointer(s.Data_))
		s.Freed_ = true
	}
}