#include "rocksdb/c.h"

// This API provides convenient C wrapper functions for rocksdb client.

/* Base */

/* Comparator */

extern rocksdb_comparator_t* api_comparator_create(uintptr_t idx);

/* Merge Operator */

extern rocksdb_mergeoperator_t* api_mergeoperator_create(uintptr_t idx);

/* CompactionFilter */

extern rocksdb_compactionfilter_t* api_compactionfilter_create(uintptr_t idx);

/* Filter Policy */

extern rocksdb_filterpolicy_t* api_filterpolicy_create(uintptr_t idx);

/* Slice Transform */

extern rocksdb_slicetransform_t* api_slicetransform_create(uintptr_t idx);