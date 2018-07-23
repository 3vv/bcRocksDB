# bcRocksDB, a Go-wrapper for RocksDB

[![Linux/Mac Build Status](https://travis-ci.org/facebook/rocksdb.svg?branch=master)](https://travis-ci.org/facebook/rocksdb)
[![Windows Build status](https://ci.appveyor.com/api/projects/status/fbgfu0so3afcno78/branch/master?svg=true)](https://ci.appveyor.com/project/Facebook/rocksdb/branch/master)
[![PPC64le Build Status](http://140.211.168.68:8080/buildStatus/icon?job=Rocksdb)](http://140.211.168.68:8080/job/Rocksdb)

## Install

You'll need to build [RocksDB](https://github.com/facebook/rocksdb) v5.5+ on your machine.

After that, you can install bcRocksDB using the following command:

    ROCKSDB="/path/to/rocksdb"
    CGO_CFLAGS="-I$ROCKSDB/include"
    CGO_LDFLAGS="-L$ROCKSDB/lib -lrocksdb -lstdc++ -lm -lz -lbz2 -lsnappy -llz4 -lzstd"
    go get github.com/3vv/bcRocksDB

Please note that this package might upgrade the required RocksDB version at any moment.
Vendoring is thus highly recommended if you require high stability.

*The [embedded CockroachDB's C-RocksDB](https://github.com/cockroachdb/c-rocksdb) is no longer supported.*