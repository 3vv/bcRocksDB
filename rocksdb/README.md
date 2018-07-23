# bcRocksDB, a Go-wrapper for RocksDB

[![Build Status](https://travis-ci.org/tecbot/leveldb.png)](https://travis-ci.org/tecbot/leveldb) [![GoDoc](https://godoc.org/github.com/tecbot/leveldb?status.png)](http://godoc.org/github.com/tecbot/leveldb)

## Install

You'll need to build [RocksDB](https://github.com/facebook/rocksdb) v5.5+ on your machine.

After that, you can install leveldb using the following command:

    ROCKSDB="/path/to/rocksdb"
    CGO_CFLAGS="-I$ROCKSDB/include"
    CGO_LDFLAGS="-L$ROCKSDB/lib -lrocksdb -lstdc++ -lm -lz -lbz2 -lsnappy -llz4 -lzstd"
    go get github.com/3vv/bcRocksDB

Please note that this package might upgrade the required RocksDB version at any moment.
Vendoring is thus highly recommended if you require high stability.

*The [embedded CockroachDB's C-RocksDB](https://github.com/cockroachdb/c-rocksdb) is no longer supported.*