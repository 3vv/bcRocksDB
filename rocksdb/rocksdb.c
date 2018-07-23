#include "rocksdb.h"
#include "_cgo_export.h"

/* Base */

void leveldb_destruct_handler(void* state) { }

/* Comparator */

rocksdb_comparator_t* leveldb_comparator_create(uintptr_t idx) {
    return rocksdb_comparator_create(
        (void*)idx,
        leveldb_destruct_handler,
        (int (*)(void*, const char*, size_t, const char*, size_t))(leveldb_comparator_compare),
        (const char *(*)(void*))(leveldb_comparator_name));
}

/* CompactionFilter */

rocksdb_compactionfilter_t* leveldb_compactionfilter_create(uintptr_t idx) {
    return rocksdb_compactionfilter_create(
        (void*)idx,
        leveldb_destruct_handler,
        (unsigned char (*)(void*, int, const char*, size_t, const char*, size_t, char**, size_t*, unsigned char*))(leveldb_compactionfilter_filter),
        (const char *(*)(void*))(leveldb_compactionfilter_name));
}

/* Filter Policy */

rocksdb_filterpolicy_t* leveldb_filterpolicy_create(uintptr_t idx) {
    return rocksdb_filterpolicy_create(
        (void*)idx,
        leveldb_destruct_handler,
        (char* (*)(void*, const char* const*, const size_t*, int, size_t*))(leveldb_filterpolicy_create_filter),
        (unsigned char (*)(void*, const char*, size_t, const char*, size_t))(leveldb_filterpolicy_key_may_match),
        leveldb_filterpolicy_delete_filter,
        (const char *(*)(void*))(leveldb_filterpolicy_name));
}

void leveldb_filterpolicy_delete_filter(void* state, const char* v, size_t s) { }

/* Merge Operator */

rocksdb_mergeoperator_t* leveldb_mergeoperator_create(uintptr_t idx) {
    return rocksdb_mergeoperator_create(
        (void*)idx,
        leveldb_destruct_handler,
        (char* (*)(void*, const char*, size_t, const char*, size_t, const char* const*, const size_t*, int, unsigned char*, size_t*))(leveldb_mergeoperator_full_merge),
        (char* (*)(void*, const char*, size_t, const char* const*, const size_t*, int, unsigned char*, size_t*))(leveldb_mergeoperator_partial_merge_multi),
        leveldb_mergeoperator_delete_value,
        (const char* (*)(void*))(leveldb_mergeoperator_name));
}

void leveldb_mergeoperator_delete_value(void* id, const char* v, size_t s) { }

/* Slice Transform */

rocksdb_slicetransform_t* leveldb_slicetransform_create(uintptr_t idx) {
    return rocksdb_slicetransform_create(
    	(void*)idx,
    	leveldb_destruct_handler,
    	(char* (*)(void*, const char*, size_t, size_t*))(leveldb_slicetransform_transform),
    	(unsigned char (*)(void*, const char*, size_t))(leveldb_slicetransform_in_domain),
    	(unsigned char (*)(void*, const char*, size_t))(leveldb_slicetransform_in_range),
    	(const char* (*)(void*))(leveldb_slicetransform_name));
}
