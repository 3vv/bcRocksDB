package opt

import (
	"../filter"
)

const (
	KiB = 1024
	MiB = KiB * 1024
	GiB = MiB * 1024
)

type Options struct {
	// OpenFilesCacheCapacity defines the capacity of the open files caching.
	// Use -1 for zero, this has same effect as specifying NoCacher to OpenFilesCacher.
	//
	// The default value is 500.
	OpenFilesCacheCapacity int

	// BlockCacheCapacity defines the capacity of the 'sorted table' block caching.
	// Use -1 for zero, this has same effect as specifying NoCacher to BlockCacher.
	//
	// The default value is 8MiB.
	BlockCacheCapacity int

	// WriteBuffer defines maximum size of a 'memdb' before flushed to
	// 'sorted table'. 'memdb' is an in-memory DB backed by an on-disk
	// unsorted journal.
	//
	// LevelDB may held up to two 'memdb' at the same time.
	//
	// The default value is 4MiB.
	WriteBuffer int

	// Filter defines an 'effective filter' to use. An 'effective filter'
	// if defined will be used to generate per-table filter block.
	// The filter name will be stored on disk.
	// During reads LevelDB will try to find matching filter from
	// 'effective filter' and 'alternative filters'.
	//
	// Filter can be changed after a DB has been created. It is recommended
	// to put old filter to the 'alternative filters' to mitigate lack of
	// filter during transition period.
	//
	// A filter is used to reduce disk reads when looking for a specific key.
	//
	// The default value is nil.
	Filter filter.Filter
}