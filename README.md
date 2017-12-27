# Go-mdbm

- Go-mdbm is a Go(golang,go-lang) bind to [Yahoo! MDBM C API](https://github.com/yahoo/mdbm).
- MDBM is a super-fast memory-mapped key/value store.
- MDBM is an ndbm work-alike hashed database library based on sdbm which is based on Per-Aake Larson’s Dynamic Hashing algorithms.
- MDBM is a high-performance, memory-mapped hash database similar to the homegrown libhash.
- The records stored in a mdbm database may have keys and values of arbitrary and variable lengths.

[![Build Status](https://travis-ci.org/torden/go-mdbm.svg?branch=master)](https://travis-ci.org/torden/go-mdbm)
[![Go Report Card](https://goreportcard.com/badge/github.com/torden/go-mdbm)](https://goreportcard.com/report/github.com/torden/go-mdbm)
[![GoDoc](https://godoc.org/github.com/torden/go-mdbm?status.svg)](https://godoc.org/github.com/torden/go-mdbm)
[![Coverage Status](https://coveralls.io/repos/github/torden/go-mdbm/badge.svg?branch=master)](https://coveralls.io/github/torden/go-mdbm?branch=master)
[![Go Walker](http://gowalker.org/api/v1/badge)](https://gowalker.org/github.com/torden/go-mdbm)
[![GitHub version](https://badge.fury.io/gh/torden%2Fgo-mdbm.svg)](https://badge.fury.io/gh/torden%2Fgo-mdbm)


## Benchmark

The following is result of Go-mdbm vs  BoltDB benchmarks for simple data storing and random fetching in them.

- [Source Code](https://github.com/torden/go-mdbm/blob/master/benchmark_test.go)
- [MDBM::Performance](http://yahoo.github.io/mdbm/guide/performance.html)

### Command

```shell
#CPU : 2 vCore
#RAM : 8G
go test -race -bench=. -run Benchmark -test.benchmem -v 
```

### Output

```
Benchmark_boltdb_Store-2                            2000           1082487 ns/op           38110 B/op         59 allocs/op
Benchmark_mdbm_Store-2                            500000              2845 ns/op              96 B/op          6 allocs/op
Benchmark_mdbm_StoreWithLock-2                    500000              2908 ns/op              96 B/op          6 allocs/op
Benchmark_boltdb_Fetch-2                          200000              9199 ns/op             496 B/op          9 allocs/op
Benchmark_mdbm_Fetch-2                           1000000              1824 ns/op              56 B/op          4 allocs/op
Benchmark_mdbm_FetchWithLock-2                   1000000              2025 ns/op              56 B/op          4 allocs/op
Benchmark_mdbm_PreLoad_Fetch-2                   1000000              1811 ns/op              56 B/op          4 allocs/op
Benchmark_mdbm_PreLoad_FetchWithLock-2            500000              2038 ns/op              56 B/op          4 allocs/op
```



## Install from Source Code

### Y! MDBM

#### Downloading

Use the master branch

```shell
git clone https://github.com/yahoo/mdbm.git
```

OR Use the release tarball, Guess you will using it

```shell
wget https://github.com/yahoo/mdbm/archive/v4.12.3.tar.gz
tar xvzf v4.12.3.tar.gz
```

#### Compiling

Refer to the https://github.com/yahoo/mdbm/blob/master/README.build
if you want to install to another path, (HIGH RECOMMEND)

```shell
cd mdbm
PREFIX=/usr/local/mdbm make install
```

## Install from Pre-build Packages

- [Ubuntu Package by Version](https://github.com/torden/go-mdbm/tree/master/pkg)
- RedHat Package by Version (as soon)
- OSX Package by Version (as soon)
- BSD Port by version (as soon)

### Ubuntu

```shell
dpkg -i pkg/ubunt/mdbm-XXXX_XXXX.dep
echo "/usr/local/mdbm/lib64/" >> /etc/ld.so.conf
```

### go-mdbm

#### Install

```
CGO_CFLAGS="-I/usr/local/mdbm/include/ -I./" \
CGO_LDFLAGS="-L/usr/local/mdbm/lib64/ -Wl,-rpath=/usr/local/mdbm/lib64/ -lmdbm" \
go get github.com/torden/go-mdbm
```

#### Download for Development or Customization

```
git clone https://github.com/torden/go-mdbm
```

##### Build

```shell
cd $GOPATH/src/github.com/torden/go-mdbm

make clean
make setup
make build
```

##### Testing

```shell
cd $GOPATH/src/github.com/torden/go-mdbm
make test
```

##### Run to Example

```shell
cd $GOPATH/src/github.com/torden/go-mdbm/example/
go run -race example.go
```

##### Miscellaneous


The following is support to development
```shell
make help
```

```shell
build:             Build the go-mdbm
clean::            Clean-up
cover:             Generate a report about coverage
coveralls::        Send a report of coverage profile to coveralls.io
help::             Show Help
installpkgs::      Install Packages
lint:              Run a LintChecker (Normal)
metalinter::       Install GoMetaLinter
pprof:             Profiling
report:            Generate the report for profiling
setup:             Setup Build Environment
strictlint:        Run a LintChecker (Strict)
test:              Run Go Test with Data Race Detection
```


## Support two compatibility branches

|*Branch*|*Support*|*test*|
|---|---|---|
|master|yes|always|
|release 4.3.x|yes|tested|

# Example

[Source Code and a sample file](https://github.com/torden/go-mdbm/tree/master/example)

## Not Support APIs

Unfortunately, the following list is not supported on now.
If you want them, please feel free to raise an issue

### Deprecated APIs

|*API*|*STATUS*|*COMMENT*|
|---|---|---|
|mdbm_save|DEPRECATED|mdbm_save is only supported for V2 MDBMs.|
|mdbm_restore|DEPRECATED|mdbm_restore is only supported for V2 MDBMs.|
|mdbm_sethash|DEPRECATED|Legacy version of mdbm_set_hash() This function has inconsistent naming, an error return value. It will be removed in a future version.|

### Only a V2 implementation

|*API*|*STATUS*|*COMMENT*|
|---|---|---|
|mdbm_stat_all_page|V3 not supported|There is only a V2 implementation. V3 not currently supported.|
|mdbm_stat_header|V3 not supported|There is only a V2 implementation. V3 not currently supported.|

### Alternative

|*API*|*COMMENT*|
|---|---|
|mdbm_fetch_buf|Fetch, FetchStr APIs|

### As soon

|*API*|*STATUS*|
|---|---|
|mdbm_cdbdump_add_record|as soon|
|mdbm_cdbdump_import|as soon|
|mdbm_dbdump_to_file|as soon|
|mdbm_dbdump_trailer_and_close|as soon|
|mdbm_dbdump_export_header|as soon|
|mdbm_set_backingstore|as soon|
|mdbm_replace_backing_store|as soon|
|mdbm_iterate|as soon|
|mdbm_prune|as soon|
|mdbm_set_stats_func|as soon|
|mdbm_chunk_iterate|as soon|


### Additional Benchmarks

#### Command
```
go test -race -bench=. -run Benchmark -test.benchmem -v 
```

##### Output

```
Benchmark_boltdb_Store-2                            2000           1072628 ns/op           37994 B/op         59 allocs/op
Benchmark_boltdb_Fetch-2                          200000              8960 ns/op             496 B/op          9 allocs/op
Benchmark_mdbm_Store-2                            500000              2923 ns/op              96 B/op          6 allocs/op
Benchmark_mdbm_StoreWithLock-2                    500000              3055 ns/op              96 B/op          6 allocs/op
Benchmark_mdbm_Fetch-2                           1000000              1823 ns/op              56 B/op          4 allocs/op
Benchmark_mdbm_FetchWithLock-2                   1000000              2093 ns/op              56 B/op          4 allocs/op
Benchmark_mdbm_PreLoad_Fetch-2                   1000000              1825 ns/op              56 B/op          4 allocs/op
Benchmark_mdbm_PreLoad_FetchWithLock-2            500000              2084 ns/op              56 B/op          4 allocs/op
```

##### DB File

|Type|File Size|Records|elapsed time|
|---|---|---|---|
|BoltDB|128K|10000|1168516ns|
|MDBM(Store)|32M|3000000|2937ns|

#### Command
```
go test -race -bench=. -run Benchmark -test.benchmem -v -test.benchtime 3s
```

##### Output

```
Benchmark_boltdb_Store-2                            5000           1168516 ns/op           39249 B/op         59 allocs/op
Benchmark_boltdb_Fetch-2                          500000              9146 ns/op             496 B/op          9 allocs/op
Benchmark_mdbm_Store-2                           2000000              2937 ns/op              96 B/op          6 allocs/op
Benchmark_mdbm_StoreWithLock-2                   1000000              3015 ns/op              96 B/op          6 allocs/op
Benchmark_mdbm_Fetch-2                           2000000              1891 ns/op              56 B/op          4 allocs/op
Benchmark_mdbm_FetchWithLock-2                   2000000              2185 ns/op              56 B/op          4 allocs/op
Benchmark_mdbm_PreLoad_Fetch-2                   2000000              1975 ns/op              56 B/op          4 allocs/op
Benchmark_mdbm_PreLoad_FetchWithLock-2           2000000              2129 ns/op              56 B/op          4 allocs/op
```

##### DB File

|Type|File Size|Records|elapsed time|
|---|---|---|---|
|BoltDB|256K|10000|1168516ns|
|MDBM(Store)|128M|3000000|2937ns|


#### Command

```
go test -race -bench=. -run Benchmark -test.benchmem -v -test.benchtime 10s
```

##### Output

```
Benchmark_boltdb_Store-2                           10000           1115691 ns/op           41109 B/op         60 allocs/op
Benchmark_mdbm_Store-2                           5000000              2933 ns/op              96 B/op          6 allocs/op
Benchmark_mdbm_StoreWithLock-2                   5000000              3098 ns/op              96 B/op          6 allocs/op
Benchmark_boltdb_Fetch-2                         2000000              9200 ns/op             496 B/op          9 allocs/op
Benchmark_mdbm_Fetch-2                          10000000              1802 ns/op              56 B/op          4 allocs/op
Benchmark_mdbm_FetchWithLock-2                  10000000              2049 ns/op              56 B/op          4 allocs/op
Benchmark_mdbm_PreLoad_Fetch-2                  10000000              1802 ns/op              56 B/op          4 allocs/op
Benchmark_mdbm_PreLoad_FetchWithLock-2          10000000              2038 ns/op              56 B/op          4 allocs/op
```

##### DB File

|Type|File Size|Records|elapsed time|
|---|---|---|---|
|BoltDB|512K|10000|1115691ns|
|MDBM(Store)|257M|5000000|2933ns|



## Links

- [Yahoo! MDBM](https://github.com/yahoo/mdbm)
- [MDBM::Concept](http://yahoo.github.io/mdbm/guide/concepts.html)
- [MDBM::Build](https://github.com/yahoo/mdbm/blob/master/README.build)
- [MDBM::Document](http://yahoo.github.io/mdbm/)
- [MDBM::FAQ](http://yahoo.github.io/mdbm/guide/faq.html)
- [DBM](https://en.wikipedia.org/wiki/Dbm)
- [BoltDB](https://github.com/boltdb/bolt)

---

*Please feel free. I hope it is helpful for you*

