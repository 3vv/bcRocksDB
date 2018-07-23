/*
Package leveldb provides the ability to create and access RocksDB databases.

leveldb.OpenDb opens and creates databases.

	bbto := leveldb.NewDefaultBlockBasedTableOptions()
	bbto.SetBlockCache(leveldb.NewLRUCache(3 << 30))
	opts := leveldb.NewDefaultOptions()
	opts.SetBlockBasedTableFactory(bbto)
	opts.SetCreateIfMissing(true)
	db, err := leveldb.OpenDb(opts, "/path/to/db")

The DB struct returned by OpenDb provides DB.Get, DB.Put, DB.Merge and DB.Delete to modify
and query the database.

	ro := leveldb.NewDefaultReadOptions()
	wo := leveldb.NewDefaultWriteOptions()
	// if ro and wo are not used again, be sure to Close them.
	err = db.Put(wo, []byte("foo"), []byte("bar"))
	...
	value, err := db.Get(ro, []byte("foo"))
	defer value.Free()
	...
	err = db.Delete(wo, []byte("foo"))

For bulk reads, use an Iterator. If you want to avoid disturbing your live
traffic while doing the bulk read, be sure to call SetFillCache(false) on the
ReadOptions you use when creating the Iterator.

	ro := leveldb.NewDefaultReadOptions()
	ro.SetFillCache(false)
	it := db.NewIterator(ro)
	defer it.Close()
	it.Seek([]byte("foo"))
	for it = it; it.Valid(); it.Next() {
		key := it.Key()
		value := it.Value()
		fmt.Printf("Key: %v Value: %v\n", key.Data(), value.Data())
		key.Free()
		value.Free()
	}
	if err := it.Err(); err != nil {
		...
	}

Batched, atomic writes can be performed with a WriteBatch and
DB.Write.

	wb := leveldb.NewWriteBatch()
	// defer wb.Close or use wb.Clear and reuse.
	wb.Delete([]byte("foo"))
	wb.Put([]byte("foo"), []byte("bar"))
	wb.Put([]byte("bar"), []byte("foo"))
	err := db.Write(wo, wb)

If your working dataset does not fit in memory, you'll want to add a bloom
filter to your database. NewBloomFilter and
BlockBasedTableOptions.SetFilterPolicy is what you want. NewBloomFilter is
amount of bits in the filter to use per key in your database.

	filter := leveldb.NewBloomFilter(10)
	bbto := leveldb.NewDefaultBlockBasedTableOptions()
	bbto.SetFilterPolicy(filter)
	opts.SetBlockBasedTableFactory(bbto)
	db, err := leveldb.OpenDb(opts, "/path/to/db")

If you're using a custom comparator in your code, be aware you may have to
make your own filter policy object.

This documentation is not a complete discussion of RocksDB. Please read the
RocksDB documentation <http://rocksdb.org/> for information on its
operation. You'll find lots of goodies there.
*/
package rocksdb