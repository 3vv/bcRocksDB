package rocksdb

//#include "api.h"
//#include <stdlib.h>
import "C"
import (
	"errors"
	"unsafe"
	"./opt"
	"./iterator"
	. "./util"
)

// DB is a reusable handle to a RocksDB database on disk, created by Open.
type DB struct {
	c    *C.rocksdb_t
	name string
	opts *Options
}

// OpenDb opens a database with the specified options.
func OpenDb(opts *Options, name string) (*DB, error) {
	var (
		cErr  *C.char
		cName = C.CString(name)
	)
	defer C.free(unsafe.Pointer(cName))
	db := C.rocksdb_open(opts.c, cName, &cErr)
	if cErr != nil {
		defer C.free(unsafe.Pointer(cErr))
		return nil, errors.New(C.GoString(cErr))
	}
	return &DB{
		name: name,
		c:    db,
		opts: opts,
	}, nil
}

func OpenFile(p string, o *opt.Options) (db *DB, err error) {
	opts := NewDefaultOptions()
  rateLimiter := NewRateLimiter(int64(o.BlockCacheCapacity), int64(o.WriteBuffer), int32(o.Filter))
  opts.SetTableCacheRemoveScanCountLimit(o.OpenFilesCacheCapacity)
	opts.SetRateLimiter(rateLimiter)
  opts.SetCreateIfMissing(true)
  opts.SetSkipLogErrorOnRecovery(false)
  return OpenDb(opts, p)
}

// OpenDbForReadOnly opens a database with the specified options for readonly usage.
func OpenDbForReadOnly(opts *Options, name string, errorIfLogFileExist bool) (*DB, error) {
	var (
		cErr  *C.char
		cName = C.CString(name)
	)
	defer C.free(unsafe.Pointer(cName))
	db := C.rocksdb_open_for_read_only(opts.c, cName, boolToChar(errorIfLogFileExist), &cErr)
	if cErr != nil {
		defer C.free(unsafe.Pointer(cErr))
		return nil, errors.New(C.GoString(cErr))
	}
	return &DB{
		name: name,
		c:    db,
		opts: opts,
	}, nil
}

// OpenDbColumnFamilies opens a database with the specified column families.
func OpenDbColumnFamilies(
	opts *Options,
	name string,
	cfNames []string,
	cfOpts []*Options,
) (*DB, []*ColumnFamilyHandle, error) {
	numColumnFamilies := len(cfNames)
	if numColumnFamilies != len(cfOpts) {
		return nil, nil, errors.New("must provide the same number of column family names and options")
	}

	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	cNames := make([]*C.char, numColumnFamilies)
	for i, s := range cfNames {
		cNames[i] = C.CString(s)
	}
	defer func() {
		for _, s := range cNames {
			C.free(unsafe.Pointer(s))
		}
	}()

	cOpts := make([]*C.rocksdb_options_t, numColumnFamilies)
	for i, o := range cfOpts {
		cOpts[i] = o.c
	}

	cHandles := make([]*C.rocksdb_column_family_handle_t, numColumnFamilies)

	var cErr *C.char
	db := C.rocksdb_open_column_families(
		opts.c,
		cName,
		C.int(numColumnFamilies),
		&cNames[0],
		&cOpts[0],
		&cHandles[0],
		&cErr,
	)
	if cErr != nil {
		defer C.free(unsafe.Pointer(cErr))
		return nil, nil, errors.New(C.GoString(cErr))
	}

	cfHandles := make([]*ColumnFamilyHandle, numColumnFamilies)
	for i, c := range cHandles {
		cfHandles[i] = NewNativeColumnFamilyHandle(c)
	}

	return &DB{
		name: name,
		c:    db,
		opts: opts,
	}, cfHandles, nil
}

// OpenDbForReadOnlyColumnFamilies opens a database with the specified column
// families in read only mode.
func OpenDbForReadOnlyColumnFamilies(
	opts *Options,
	name string,
	cfNames []string,
	cfOpts []*Options,
	errorIfLogFileExist bool,
) (*DB, []*ColumnFamilyHandle, error) {
	numColumnFamilies := len(cfNames)
	if numColumnFamilies != len(cfOpts) {
		return nil, nil, errors.New("must provide the same number of column family names and options")
	}

	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	cNames := make([]*C.char, numColumnFamilies)
	for i, s := range cfNames {
		cNames[i] = C.CString(s)
	}
	defer func() {
		for _, s := range cNames {
			C.free(unsafe.Pointer(s))
		}
	}()

	cOpts := make([]*C.rocksdb_options_t, numColumnFamilies)
	for i, o := range cfOpts {
		cOpts[i] = o.c
	}

	cHandles := make([]*C.rocksdb_column_family_handle_t, numColumnFamilies)

	var cErr *C.char
	db := C.rocksdb_open_for_read_only_column_families(
		opts.c,
		cName,
		C.int(numColumnFamilies),
		&cNames[0],
		&cOpts[0],
		&cHandles[0],
		boolToChar(errorIfLogFileExist),
		&cErr,
	)
	if cErr != nil {
		defer C.free(unsafe.Pointer(cErr))
		return nil, nil, errors.New(C.GoString(cErr))
	}

	cfHandles := make([]*ColumnFamilyHandle, numColumnFamilies)
	for i, c := range cHandles {
		cfHandles[i] = NewNativeColumnFamilyHandle(c)
	}

	return &DB{
		name: name,
		c:    db,
		opts: opts,
	}, cfHandles, nil
}

// ListColumnFamilies lists the names of the column families in the DB.
func ListColumnFamilies(opts *Options, name string) ([]string, error) {
	var (
		cErr  *C.char
		cLen  C.size_t
		cName = C.CString(name)
	)
	defer C.free(unsafe.Pointer(cName))
	cNames := C.rocksdb_list_column_families(opts.c, cName, &cLen, &cErr)
	if cErr != nil {
		defer C.free(unsafe.Pointer(cErr))
		return nil, errors.New(C.GoString(cErr))
	}
	namesLen := int(cLen)
	names := make([]string, namesLen)
	cNamesArr := (*[1 << 30]*C.char)(unsafe.Pointer(cNames))[:namesLen:namesLen]
	for i, n := range cNamesArr {
		names[i] = C.GoString(n)
	}
	C.rocksdb_list_column_families_destroy(cNames, cLen)
	return names, nil
}

// UnsafeGetDB returns the underlying c rocksdb instance.
func (db *DB) UnsafeGetDB() unsafe.Pointer {
	return unsafe.Pointer(db.c)
}

// Name returns the name of the database.
func (db *DB) Name() string {
	return db.name
}

// Get returns the data associated with the key from the database.
func (db *DB) Get(/*opts *ReadOptions, */key []byte, opts *ReadOptions) (/**Slice*/[]byte, error) {/*
	var (
		cErr    *C.char
		cValLen C.size_t
		cKey    = ByteToChar(key)
	)
	cValue := C.rocksdb_get(db.c, opts.c, cKey, C.size_t(len(key)), &cValLen, &cErr)
	if cErr != nil {
		defer C.free(unsafe.Pointer(cErr))
		return nil, errors.New(C.GoString(cErr))
	}
	return NewSlice(cValue, cValLen), nil*/
	return db.GetBytes(opts, key)
}

// GetBytes is like Get but returns a copy of the data.
func (db *DB) GetBytes(opts *ReadOptions, key []byte) ([]byte, error) {
	var (
		cErr    *C.char
		cValLen C.size_t
		cKey    = /*ByteToChar(key)*/(*C.char)(unsafe.Pointer(&key[0]))
	)
	cValue := C.rocksdb_get(db.c, opts.c, cKey, C.size_t(len(key)), &cValLen, &cErr)
	if cErr != nil {
		defer C.free(unsafe.Pointer(cErr))
		return nil, errors.New(C.GoString(cErr))
	}
	if cValue == nil {
		return nil, nil
	}
	defer C.free(unsafe.Pointer(cValue))
	return C.GoBytes(unsafe.Pointer(cValue), C.int(cValLen)), nil
}

// GetCF returns the data associated with the key from the database and column family.
func (db *DB) GetCF(opts *ReadOptions, cf *ColumnFamilyHandle, key []byte) (*Slice, error) {
	var (
		cErr    *C.char
		cValLen C.size_t
		cKey    = /*ByteToChar(key)*/(*C.char)(unsafe.Pointer(&key[0]))
	)
	cValue := C.rocksdb_get_cf(db.c, opts.c, cf.c, cKey, C.size_t(len(key)), &cValLen, &cErr)
	if cErr != nil {
		defer C.free(unsafe.Pointer(cErr))
		return nil, errors.New(C.GoString(cErr))
	}
	return /*NewSlice(cValue, cValLen), nil*/StringToSlice(C.GoStringN(cValue, (C.int)(cValLen))), nil
}

// Put writes data associated with a key to the database.
func (db *DB) Put(/*opts *WriteOptions, */key, value []byte, opts *WriteOptions) error {
	var (
		cErr   *C.char
		cKey   = /*ByteToChar(key)*/(*C.char)(unsafe.Pointer(&key[0]))
		cValue = /*ByteToChar(value)*/(*C.char)(unsafe.Pointer(&value[0]))
	)
	C.rocksdb_put(db.c, opts.c, cKey, C.size_t(len(key)), cValue, C.size_t(len(value)), &cErr)
	if cErr != nil {
		defer C.free(unsafe.Pointer(cErr))
		return errors.New(C.GoString(cErr))
	}
	return nil
}

// PutCF writes data associated with a key to the database and column family.
func (db *DB) PutCF(opts *WriteOptions, cf *ColumnFamilyHandle, key, value []byte) error {
	var (
		cErr   *C.char
		cKey   = /*ByteToChar(key)*/(*C.char)(unsafe.Pointer(&key[0]))
		cValue = /*ByteToChar(value)*/(*C.char)(unsafe.Pointer(&value[0]))
	)
	C.rocksdb_put_cf(db.c, opts.c, cf.c, cKey, C.size_t(len(key)), cValue, C.size_t(len(value)), &cErr)
	if cErr != nil {
		defer C.free(unsafe.Pointer(cErr))
		return errors.New(C.GoString(cErr))
	}
	return nil
}

// Delete removes the data associated with the key from the database.
func (db *DB) Delete(/*opts *WriteOptions, */key []byte, opts *WriteOptions) error {
	var (
		cErr *C.char
		cKey = /*ByteToChar(key)*/(*C.char)(unsafe.Pointer(&key[0]))
	)
	C.rocksdb_delete(db.c, opts.c, cKey, C.size_t(len(key)), &cErr)
	if cErr != nil {
		defer C.free(unsafe.Pointer(cErr))
		return errors.New(C.GoString(cErr))
	}
	return nil
}

// DeleteCF removes the data associated with the key from the database and column family.
func (db *DB) DeleteCF(opts *WriteOptions, cf *ColumnFamilyHandle, key []byte) error {
	var (
		cErr *C.char
		cKey = /*ByteToChar(key)*/(*C.char)(unsafe.Pointer(&key[0]))
	)
	C.rocksdb_delete_cf(db.c, opts.c, cf.c, cKey, C.size_t(len(key)), &cErr)
	if cErr != nil {
		defer C.free(unsafe.Pointer(cErr))
		return errors.New(C.GoString(cErr))
	}
	return nil
}

// Merge merges the data associated with the key with the actual data in the database.
func (db *DB) Merge(opts *WriteOptions, key []byte, value []byte) error {
	var (
		cErr   *C.char
		cKey   = /*ByteToChar(key)*/(*C.char)(unsafe.Pointer(&key[0]))
		cValue = /*ByteToChar(value)*/(*C.char)(unsafe.Pointer(&value[0]))
	)
	C.rocksdb_merge(db.c, opts.c, cKey, C.size_t(len(key)), cValue, C.size_t(len(value)), &cErr)
	if cErr != nil {
		defer C.free(unsafe.Pointer(cErr))
		return errors.New(C.GoString(cErr))
	}
	return nil
}

// MergeCF merges the data associated with the key with the actual data in the
// database and column family.
func (db *DB) MergeCF(opts *WriteOptions, cf *ColumnFamilyHandle, key []byte, value []byte) error {
	var (
		cErr   *C.char
		cKey   = /*ByteToChar(key)*/(*C.char)(unsafe.Pointer(&key[0]))
		cValue = /*ByteToChar(value)*/(*C.char)(unsafe.Pointer(&value[0]))
	)
	C.rocksdb_merge_cf(db.c, opts.c, cf.c, cKey, C.size_t(len(key)), cValue, C.size_t(len(value)), &cErr)
	if cErr != nil {
		defer C.free(unsafe.Pointer(cErr))
		return errors.New(C.GoString(cErr))
	}
	return nil
}

// Write writes a WriteBatch to the database
func (db *DB) Write(/*opts *WriteOptions, batch *WriteBatch*/batch *Batch, opts *WriteOptions) error {
	var cErr *C.char
	C.rocksdb_write(db.c, opts.c, batch.c, &cErr)
	if cErr != nil {
		defer C.free(unsafe.Pointer(cErr))
		return errors.New(C.GoString(cErr))
	}
	return nil
}

// NewIterator returns an Iterator over the the database that uses the
// ReadOptions given.
func (db *DB) NewIterator(slice *Range, opts *ReadOptions) *iterator.Iterator {
	cIter := C.rocksdb_create_iterator(db.c, opts.c)
	return iterator.NewNativeIterator(unsafe.Pointer(cIter))
}

// NewIteratorCF returns an Iterator over the the database and column family
// that uses the ReadOptions given.
func (db *DB) NewIteratorCF(opts *ReadOptions, cf *ColumnFamilyHandle) *iterator.Iterator {
	cIter := C.rocksdb_create_iterator_cf(db.c, opts.c, cf.c)
	return iterator.NewNativeIterator(unsafe.Pointer(cIter))
}

// NewSnapshot creates a new snapshot of the database.
func (db *DB) NewSnapshot() *Snapshot {
	cSnap := C.rocksdb_create_snapshot(db.c)
	return NewNativeSnapshot(cSnap)
}

// ReleaseSnapshot releases the snapshot and its resources.
func (db *DB) ReleaseSnapshot(snapshot *Snapshot) {
	C.rocksdb_release_snapshot(db.c, snapshot.c)
	snapshot.c = nil
}

// GetProperty returns the value of a database property.
func (db *DB) GetProperty(propName string) (string, error) {
	cprop := C.CString(propName)
	defer C.free(unsafe.Pointer(cprop))
	cValue := C.rocksdb_property_value(db.c, cprop)
	defer C.free(unsafe.Pointer(cValue))
	return C.GoString(cValue), nil
}

// GetPropertyCF returns the value of a database property.
func (db *DB) GetPropertyCF(propName string, cf *ColumnFamilyHandle) string {
	cProp := C.CString(propName)
	defer C.free(unsafe.Pointer(cProp))
	cValue := C.rocksdb_property_value_cf(db.c, cf.c, cProp)
	defer C.free(unsafe.Pointer(cValue))
	return C.GoString(cValue)
}

// CreateColumnFamily create a new column family.
func (db *DB) CreateColumnFamily(opts *Options, name string) (*ColumnFamilyHandle, error) {
	var (
		cErr  *C.char
		cName = C.CString(name)
	)
	defer C.free(unsafe.Pointer(cName))
	cHandle := C.rocksdb_create_column_family(db.c, opts.c, cName, &cErr)
	if cErr != nil {
		defer C.free(unsafe.Pointer(cErr))
		return nil, errors.New(C.GoString(cErr))
	}
	return NewNativeColumnFamilyHandle(cHandle), nil
}

// DropColumnFamily drops a column family.
func (db *DB) DropColumnFamily(c *ColumnFamilyHandle) error {
	var cErr *C.char
	C.rocksdb_drop_column_family(db.c, c.c, &cErr)
	if cErr != nil {
		defer C.free(unsafe.Pointer(cErr))
		return errors.New(C.GoString(cErr))
	}
	return nil
}

// GetApproximateSizes returns the approximate number of bytes of file system
// space used by one or more key ranges.
//
// The keys counted will begin at Range.Start and end on the key before
// Range.Limit.
func (db *DB) GetApproximateSizes(ranges []Range) []uint64 {
	sizes := make([]uint64, len(ranges))
	if len(ranges) == 0 {
		return sizes
	}

	cStarts := make([]*C.char, len(ranges))
	cLimits := make([]*C.char, len(ranges))
	cStartLens := make([]C.size_t, len(ranges))
	cLimitLens := make([]C.size_t, len(ranges))
	for i, r := range ranges {
		cStarts[i] = /*ByteToChar(r.Start)*/(*C.char)(unsafe.Pointer(&r.Start[0]))
		cStartLens[i] = C.size_t(len(r.Start))
		cLimits[i] = /*ByteToChar(r.Limit)*/(*C.char)(unsafe.Pointer(&r.Limit[0]))
		cLimitLens[i] = C.size_t(len(r.Limit))
	}

	C.rocksdb_approximate_sizes(
		db.c,
		C.int(len(ranges)),
		&cStarts[0],
		&cStartLens[0],
		&cLimits[0],
		&cLimitLens[0],
		(*C.uint64_t)(&sizes[0]))

	return sizes
}

// GetApproximateSizesCF returns the approximate number of bytes of file system
// space used by one or more key ranges in the column family.
//
// The keys counted will begin at Range.Start and end on the key before
// Range.Limit.
func (db *DB) GetApproximateSizesCF(cf *ColumnFamilyHandle, ranges []Range) []uint64 {
	sizes := make([]uint64, len(ranges))
	if len(ranges) == 0 {
		return sizes
	}

	cStarts := make([]*C.char, len(ranges))
	cLimits := make([]*C.char, len(ranges))
	cStartLens := make([]C.size_t, len(ranges))
	cLimitLens := make([]C.size_t, len(ranges))
	for i, r := range ranges {
		cStarts[i] = /*ByteToChar(r.Start)*/(*C.char)(unsafe.Pointer(&r.Start[0]))
		cStartLens[i] = C.size_t(len(r.Start))
		cLimits[i] = /*ByteToChar(r.Limit)*/(*C.char)(unsafe.Pointer(&r.Limit[0]))
		cLimitLens[i] = C.size_t(len(r.Limit))
	}

	C.rocksdb_approximate_sizes_cf(
		db.c,
		cf.c,
		C.int(len(ranges)),
		&cStarts[0],
		&cStartLens[0],
		&cLimits[0],
		&cLimitLens[0],
		(*C.uint64_t)(&sizes[0]))

	return sizes
}

// LiveFileMetadata is a metadata which is associated with each SST file.
type LiveFileMetadata struct {
	Name        string
	Level       int
	Size        int64
	SmallestKey []byte
	LargestKey  []byte
}

// GetLiveFilesMetaData returns a list of all table files with their
// level, start key and end key.
func (db *DB) GetLiveFilesMetaData() []LiveFileMetadata {
	lf := C.rocksdb_livefiles(db.c)
	defer C.rocksdb_livefiles_destroy(lf)

	count := C.rocksdb_livefiles_count(lf)
	liveFiles := make([]LiveFileMetadata, int(count))
	for i := C.int(0); i < count; i++ {
		var liveFile LiveFileMetadata
		liveFile.Name = C.GoString(C.rocksdb_livefiles_name(lf, i))
		liveFile.Level = int(C.rocksdb_livefiles_level(lf, i))
		liveFile.Size = int64(C.rocksdb_livefiles_size(lf, i))

		var cSize C.size_t
		key := C.rocksdb_livefiles_smallestkey(lf, i, &cSize)
		liveFile.SmallestKey = C.GoBytes(unsafe.Pointer(key), C.int(cSize))

		key = C.rocksdb_livefiles_largestkey(lf, i, &cSize)
		liveFile.LargestKey = C.GoBytes(unsafe.Pointer(key), C.int(cSize))
		liveFiles[int(i)] = liveFile
	}
	return liveFiles
}

// CompactRange runs a manual compaction on the Range of keys given. This is
// not likely to be needed for typical usage.
func (db *DB) CompactRange(r Range) error {
	cStart := /*ByteToChar(r.Start)*/(*C.char)(unsafe.Pointer(&r.Start[0]))
	cLimit := /*ByteToChar(r.Limit)*/(*C.char)(unsafe.Pointer(&r.Limit[0]))
  C.rocksdb_compact_range(db.c, cStart, C.size_t(len(r.Start)), cLimit, C.size_t(len(r.Limit)))
  return nil
}

// CompactRangeCF runs a manual compaction on the Range of keys given on the
// given column family. This is not likely to be needed for typical usage.
func (db *DB) CompactRangeCF(cf *ColumnFamilyHandle, r Range) {
	cStart := /*ByteToChar(r.Start)*/(*C.char)(unsafe.Pointer(&r.Start[0]))
	cLimit := /*ByteToChar(r.Limit)*/(*C.char)(unsafe.Pointer(&r.Limit[0]))
	C.rocksdb_compact_range_cf(db.c, cf.c, cStart, C.size_t(len(r.Start)), cLimit, C.size_t(len(r.Limit)))
}

// Flush triggers a manuel flush for the database.
func (db *DB) Flush(opts *FlushOptions) error {
	var cErr *C.char
	C.rocksdb_flush(db.c, opts.c, &cErr)
	if cErr != nil {
		defer C.free(unsafe.Pointer(cErr))
		return errors.New(C.GoString(cErr))
	}
	return nil
}

// DisableFileDeletions disables file deletions and should be used when backup the database.
func (db *DB) DisableFileDeletions() error {
	var cErr *C.char
	C.rocksdb_disable_file_deletions(db.c, &cErr)
	if cErr != nil {
		defer C.free(unsafe.Pointer(cErr))
		return errors.New(C.GoString(cErr))
	}
	return nil
}

// EnableFileDeletions enables file deletions for the database.
func (db *DB) EnableFileDeletions(force bool) error {
	var cErr *C.char
	C.rocksdb_enable_file_deletions(db.c, boolToChar(force), &cErr)
	if cErr != nil {
		defer C.free(unsafe.Pointer(cErr))
		return errors.New(C.GoString(cErr))
	}
	return nil
}

// DeleteFile deletes the file name from the db directory and update the internal state to
// reflect that. Supports deletion of sst and log files only. 'name' must be
// path relative to the db directory. eg. 000001.sst, /archive/000003.log.
func (db *DB) DeleteFile(name string) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	C.rocksdb_delete_file(db.c, cName)
}

// IngestExternalFile loads a list of external SST files.
func (db *DB) IngestExternalFile(filePaths []string, opts *IngestExternalFileOptions) error {
	cFilePaths := make([]*C.char, len(filePaths))
	for i, s := range filePaths {
		cFilePaths[i] = C.CString(s)
	}
	defer func() {
		for _, s := range cFilePaths {
			C.free(unsafe.Pointer(s))
		}
	}()

	var cErr *C.char

	C.rocksdb_ingest_external_file(
		db.c,
		&cFilePaths[0],
		C.size_t(len(filePaths)),
		opts.c,
		&cErr,
	)

	if cErr != nil {
		defer C.free(unsafe.Pointer(cErr))
		return errors.New(C.GoString(cErr))
	}
	return nil
}

// IngestExternalFileCF loads a list of external SST files for a column family.
func (db *DB) IngestExternalFileCF(handle *ColumnFamilyHandle, filePaths []string, opts *IngestExternalFileOptions) error {
	cFilePaths := make([]*C.char, len(filePaths))
	for i, s := range filePaths {
		cFilePaths[i] = C.CString(s)
	}
	defer func() {
		for _, s := range cFilePaths {
			C.free(unsafe.Pointer(s))
		}
	}()

	var cErr *C.char

	C.rocksdb_ingest_external_file_cf(
		db.c,
		handle.c,
		&cFilePaths[0],
		C.size_t(len(filePaths)),
		opts.c,
		&cErr,
	)

	if cErr != nil {
		defer C.free(unsafe.Pointer(cErr))
		return errors.New(C.GoString(cErr))
	}
	return nil
}

// NewCheckpoint creates a new Checkpoint for this db.
func (db *DB) NewCheckpoint() (*Checkpoint, error) {
	var (
		cErr *C.char
	)
	cCheckpoint := C.rocksdb_checkpoint_object_create(
		db.c, &cErr,
	)
	if cErr != nil {
		defer C.free(unsafe.Pointer(cErr))
		return nil, errors.New(C.GoString(cErr))
	}

	return NewNativeCheckpoint(cCheckpoint), nil
}

// Close closes the database.
func (db *DB) Close() error {
  C.rocksdb_close(db.c)
  return nil
}

// DestroyDb removes a database entirely, removing everything from the
// filesystem.
func DestroyDb(name string, opts *Options) error {
	var (
		cErr  *C.char
		cName = C.CString(name)
	)
	defer C.free(unsafe.Pointer(cName))
	C.rocksdb_destroy_db(opts.c, cName, &cErr)
	if cErr != nil {
		defer C.free(unsafe.Pointer(cErr))
		return errors.New(C.GoString(cErr))
	}
	return nil
}

// RepairDb repairs a database.
func RepairDb(name string, opts *Options) error {
	var (
		cErr  *C.char
		cName = C.CString(name)
	)
	defer C.free(unsafe.Pointer(cName))
	C.rocksdb_repair_db(opts.c, cName, &cErr)
	if cErr != nil {
		defer C.free(unsafe.Pointer(cErr))
		return errors.New(C.GoString(cErr))
	}
	return nil
}