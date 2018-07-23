#include "api.h"
#include "_cgo_export.h"

/* Base */

void api_destruct_handler(void *state) {}

/* Comparator */

rocksdb_comparator_t *api_comparator_create(uintptr_t idx) {
  return rocksdb_comparator_create(
    (void *)idx,
    api_destruct_handler,
    (int (*)(void *, const char *, size_t, const char *, size_t))(itf_comparator_compare),
    (const char *(*)(void *))(itf_comparator_name));
}

/* Merge Operator */

void itf_mergeoperator_delete_value(void *id, const char *v, size_t s) {}

rocksdb_mergeoperator_t *api_mergeoperator_create(uintptr_t idx) {
  return rocksdb_mergeoperator_create(
    (void *)idx,
    api_destruct_handler,
    (char *(*)(void *, const char *, size_t, const char *, size_t, const char *const *, const size_t *, int, unsigned char *, size_t *))(itf_mergeoperator_full_merge),
    (char *(*)(void *, const char *, size_t, const char *const *, const size_t *, int, unsigned char *, size_t *))(itf_mergeoperator_partial_merge_multi),
    itf_mergeoperator_delete_value,
    (const char *(*)(void *))(itf_mergeoperator_name));
}

/* CompactionFilter */

rocksdb_compactionfilter_t *api_compactionfilter_create(uintptr_t idx) {
  return rocksdb_compactionfilter_create(
    (void *)idx,
    api_destruct_handler,
    (unsigned char (*)(void *, int, const char *, size_t, const char *, size_t, char **, size_t *, unsigned char *))(itf_compactionfilter_filter),
    (const char *(*)(void *))(itf_compactionfilter_name));
}

/* Filter Policy */

void itf_filterpolicy_delete_filter(void *state, const char *v, size_t s) {}

rocksdb_filterpolicy_t *api_filterpolicy_create(uintptr_t idx) {
  return rocksdb_filterpolicy_create(
    (void *)idx,
    api_destruct_handler,
    (char *(*)(void *, const char *const *, const size_t *, int, size_t *))(itf_filterpolicy_create_filter),
    (unsigned char (*)(void *, const char *, size_t, const char *, size_t))(itf_filterpolicy_key_may_match),
    itf_filterpolicy_delete_filter,
    (const char *(*)(void *))(itf_filterpolicy_name));
}

/* Slice Transform */

rocksdb_slicetransform_t *api_slicetransform_create(uintptr_t idx) {
  return rocksdb_slicetransform_create(
    (void *)idx,
    api_destruct_handler,
    (char *(*)(void *, const char *, size_t, size_t *))(itf_slicetransform_transform),
    (unsigned char (*)(void *, const char *, size_t))(itf_slicetransform_in_domain),
    (unsigned char (*)(void *, const char *, size_t))(itf_slicetransform_in_range),
    (const char *(*)(void *))(itf_slicetransform_name));
}