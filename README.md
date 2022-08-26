
# Go-mdbm

- Go-mdbm is a Go(golang,go-lang) binds to [Yahoo! MDBM C API](https://github.com/yahoo/mdbm).
- MDBM is a super-fast memory-mapped key/value store.
- MDBM is an ndbm work-alike hashed database library based on sdbm which is based on Per-Aake Larson’s Dynamic Hashing algorithms.
- MDBM is a high-performance, memory-mapped hash database similar to the homegrown libhash.
- The records stored in a mdbm database may have keys and values of arbitrary and variable lengths.

[![Build Status](https://github.com/torden/go-mdbm/actions/workflows/go.yml/badge.svg)](https://github.com/torden/go-mdbm/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/torden/go-mdbm)](https://goreportcard.com/report/github.com/torden/go-mdbm)
[![GoDoc](https://godoc.org/github.com/torden/go-mdbm?status.svg)](https://godoc.org/github.com/torden/go-mdbm)
[![codecov](https://codecov.io/gh/torden/go-mdbm/branch/master/graph/badge.svg?token=Cb94xJMoYW)](https://codecov.io/gh/torden/go-mdbm)
[![Coverage Status](https://coveralls.io/repos/github/torden/go-mdbm/badge.svg?branch=master)](https://coveralls.io/github/torden/go-mdbm?branch=master)
[![GitHub version](https://img.shields.io/github/v/release/torden/go-mdbm)](https://github.com/torden/go-mdbm)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)


## Table of Contents


- [Install from Source Code](#install-from-source-code)
    - [Y! MDBM](#y!-mdbm)
    - [Downloading](#downloading)
    - [Compiling](#compiling)
- [Install from Pre-build Packages](#install-from-pre-build-packages)
    - [Ubuntu](#ubuntu)
    - [go-mdbm](#go-mdbm)
    - [Install](#install)
- [Download for Development or Customization](#download-for-development-or-customization)
    - [Build](#build)
    - [Testing](#testing)
    - [Run to Example](#run-to-example)
    - [Miscellaneous](#miscellaneous)
    - [Support two compatibility branches](#support-two-compatibility-branches)
- [Not Support APIs](#not-support-apis)
    - [Deprecated APIs](#deprecated-apis)
    - [Only a V2 implementation](#only-a-v2-implementation)
    - [Alternative](#alternative)
    - [As soon](#as-soon)
- [Examples](#examples)
- [Benchmark](#benchmark)
    - [Spec](#spec)
    - [Command](#command)
    - [Output](#output)
- [Additional Benchmarks](#additional-benchmarks)
    - [Spec](#spec)
        - [Command](#command)
        - [Output](#output)
        - [DB File](#db-file)
    - [Command](#command)
        - [Output](#output)
        - [DB File](#db-file)
    - [Command](#command)
        - [Output](#output)
        - [DB File](#db-file)
    - [Spec](#spec)
    - [Command](#command)
        - [Output](#output)
        - [DB File](#db-file)
    - [Command](#command)
        - [Output](#output)
        - [DB File](#db-file)
    - [Command](#command)
        - [Output](#output)
        - [DB File](#db-file)
- [Links](#links)


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

- [Ubuntu Package](https://github.com/torden/go-mdbm/tree/master/pkg)
- [RedHat(CentOS) Package](https://github.com/torden/go-mdbm/tree/master/pkg)
- OSX Package (as soon)
- BSD Port (as soon)

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
|mdbm_set_backingstore|as soon|
|mdbm_cdbdump_add_record|as soon|
|mdbm_cdbdump_import|as soon|
|mdbm_dbdump_to_file|as soon|
|mdbm_dbdump_trailer_and_close|as soon|
|mdbm_dbdump_export_header|as soon|
|mdbm_iterate|as soon|
|mdbm_prune|as soon|
|mdbm_set_stats_func|as soon|
|mdbm_chunk_iterate|as soon|

## Examples

See the [Documentations](https://github.com/torden/go-mdbm/tree/master/example) for more details


## Benchmark

The following is result of Go-mdbm vs  BoltDB benchmarks for simple data storing and random fetching in them.

- [Source Code](https://github.com/torden/go-mdbm/blob/master/benchmark_test.go)
- [MDBM::Performance](http://yahoo.github.io/mdbm/guide/performance.html)

### Spec

|Type|Spec|
|---|---|
|Machine|VM(VirtualBox)|
|OS|Ubuntu 17.10 (Artful Aardvark)|
|CPU|2 vCore|
|RAM|8G|
|BoltDB Ver.|9da3174 on 20 Nov|
|Mdbm Ver.|893f7a8 on 26 Jul|

### Command

```shell
CGO_CFLAGS="-I/usr/local/mdbm/include/ -I./" CGO_LDFLAGS="-L/usr/local/mdbm/lib64/ -Wl,-rpath=/usr/local/mdbm/lib64/ -lmdbm" \
go test -race -bench=. -run Benchmark -test.benchmem -v
```

### Output

|func|count of loop|nano-seconds per loop|bytes per operation|allocations per operation|
|---|--:|--:|--:|--:|
|Benchmark_boltdb_Store-2                  |          2000      |     1082487 ns/op    |       38110 B/op     |    59 allocs/op|
|Benchmark_mdbm_Store-2                    |        500000      |        2845 ns/op    |          96 B/op     |     6 allocs/op|
|Benchmark_mdbm_StoreWithLock-2            |        500000      |        2908 ns/op    |          96 B/op     |     6 allocs/op|
|Benchmark_boltdb_Fetch-2                  |        200000      |        9199 ns/op    |         496 B/op     |     9 allocs/op|
|Benchmark_mdbm_Fetch-2                    |       1000000      |        1824 ns/op    |          56 B/op     |     4 allocs/op|
|Benchmark_mdbm_FetchWithLock-2            |       1000000      |        2025 ns/op    |          56 B/op     |     4 allocs/op|
|Benchmark_mdbm_PreLoad_Fetch-2            |       1000000      |        1811 ns/op    |          56 B/op     |     4 allocs/op|
|Benchmark_mdbm_PreLoad_FetchWithLock-2    |        500000      |        2038 ns/op    |          56 B/op     |     4 allocs/op|

### Additional Benchmarks

#### Spec

|Type|Spec|
|---|---|
|Machine|VM(VirtualBox)|
|OS|Ubuntu 17.10 (Artful Aardvark)|
|CPU|2 vCore|
|RAM|8G|
|BoltDB Ver.|9da3174 on 20 Nov|
|Mdbm Ver.|893f7a8 on 26 Jul|


#### Command
```
CGO_CFLAGS="-I/usr/local/mdbm/include/ -I./" CGO_LDFLAGS="-L/usr/local/mdbm/lib64/ -Wl,-rpath=/usr/local/mdbm/lib64/ -lmdbm"
go test -race -bench=. -run Benchmark -test.benchmem -v \
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

|Type|File Size|Times|elapsed time|
|---|---|---|---|
|BoltDB|128K|10000|1168516ns|
|MDBM(Store)|32M|3000000|2937ns|

#### Command
```
CGO_CFLAGS="-I/usr/local/mdbm/include/ -I./" CGO_LDFLAGS="-L/usr/local/mdbm/lib64/ -Wl,-rpath=/usr/local/mdbm/lib64/ -lmdbm" \
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

|Type|File Size|Times|elapsed time|
|---|---|---|---|
|BoltDB|256K|10000|1168516ns|
|MDBM(Store)|128M|3000000|2937ns|


#### Command

```
CGO_CFLAGS="-I/usr/local/mdbm/include/ -I./" CGO_LDFLAGS="-L/usr/local/mdbm/lib64/ -Wl,-rpath=/usr/local/mdbm/lib64/ -lmdbm" \
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

|Type|File Size|Times|elapsed time|
|---|---|---|---|
|BoltDB|512K|10000|1115691ns|
|MDBM(Store)|257M|5000000|2933ns|


#### Spec

|Type|Spec|
|---|---|
|Machine|Physical|
|OS|Ubuntu 17.10 (Artful Aardvark)|
|CPU|8 Core (Intel i7)|
|RAM|16G|
|HDD|SSD|
|BoltDB Ver.|9da3174 on 20 Nov|
|Mdbm Ver.|893f7a8 on 26 Jul|

#### Command

```
#CPU : 8core
#RAM : 16g
CGO_CFLAGS="-I/usr/local/mdbm/include/ -I./" CGO_LDFLAGS="-L/usr/local/mdbm/lib64/ -Wl,-rpath=/usr/local/mdbm/lib64/ -lmdbm" \
go test -race -bench=. -run Benchmark -test.benchmem -v
```

##### Output

```
Benchmark_boltdb_Store-8                             300           6138312 ns/op           32704 B/op         55 allocs/op
Benchmark_mdbm_Store-8                            200000              5235 ns/op              96 B/op          6 allocs/op
Benchmark_mdbm_StoreWithLock-8                    200000              5749 ns/op              96 B/op          6 allocs/op
Benchmark_boltdb_Fetch-8                          100000             15322 ns/op             496 B/op          9 allocs/op
Benchmark_mdbm_Fetch-8                            500000              2852 ns/op              56 B/op          4 allocs/op
Benchmark_mdbm_FetchWithLock-8                    300000              3713 ns/op              56 B/op          4 allocs/op
Benchmark_mdbm_PreLoad_Fetch-8                    500000              2829 ns/op              56 B/op          4 allocs/op
Benchmark_mdbm_PreLoad_FetchWithLock-8            500000              3436 ns/op              56 B/op          4 allocs/op
```

##### DB File

|Type|File Size|Times|elapsed time|
|---|---|---|---|
|BoltDB|64K|300|6138312ns|
|MDBM(Store)|16M|200000|5235ns|


#### Command

```
CGO_CFLAGS="-I/usr/local/mdbm/include/ -I./" CGO_LDFLAGS="-L/usr/local/mdbm/lib64/ -Wl,-rpath=/usr/local/mdbm/lib64/ -lmdbm" \
go test -race -bench=. -run Benchmark -test.benchmem -v -test.benchtime 3s
```

##### Output

```
Benchmark_boltdb_Store-8                            1000           6283664 ns/op           37533 B/op         58 allocs/op
Benchmark_mdbm_Store-8                           1000000              4780 ns/op              96 B/op          6 allocs/op
Benchmark_mdbm_StoreWithLock-8                   1000000              5360 ns/op              96 B/op          6 allocs/op
Benchmark_boltdb_Fetch-8                          300000             14556 ns/op             496 B/op          9 allocs/op
Benchmark_mdbm_Fetch-8                           2000000              2772 ns/op              56 B/op          4 allocs/op
Benchmark_mdbm_FetchWithLock-8                   1000000              3104 ns/op              56 B/op          4 allocs/op
Benchmark_mdbm_PreLoad_Fetch-8                   2000000              2527 ns/op              56 B/op          4 allocs/op
Benchmark_mdbm_PreLoad_FetchWithLock-8           1000000              3256 ns/op              56 B/op          4 allocs/op
```

##### DB File

|Type|File Size|Times|elapsed time|
|---|---|---|---|
|BoltDB|128K|1000|6283664ns|
|MDBM(Store)|64M|1000000|4780ns|

#### Command

```
CGO_CFLAGS="-I/usr/local/mdbm/include/ -I./" CGO_LDFLAGS="-L/usr/local/mdbm/lib64/ -Wl,-rpath=/usr/local/mdbm/lib64/ -lmdbm" \
go test -race -bench=. -run Benchmark -test.benchmem -v -test.benchtime 10s
```

##### Output

```
Benchmark_boltdb_Store-8                            2000           6133872 ns/op           28366 B/op         59 allocs/op
Benchmark_mdbm_Store-8                           3000000              5377 ns/op              96 B/op          6 allocs/op
Benchmark_mdbm_StoreWithLock-8                   3000000              5145 ns/op              96 B/op          6 allocs/op
Benchmark_boltdb_Fetch-8                         1000000             15703 ns/op             496 B/op          9 allocs/op
Benchmark_mdbm_Fetch-8                           5000000              2631 ns/op              56 B/op          4 allocs/op
Benchmark_mdbm_FetchWithLock-8                   5000000              3245 ns/op              56 B/op          4 allocs/op
Benchmark_mdbm_PreLoad_Fetch-8                   5000000              2598 ns/op              56 B/op          4 allocs/op
Benchmark_mdbm_PreLoad_FetchWithLock-8           5000000              3379 ns/op              56 B/op          4 allocs/op
```

##### DB File

|Type|File Size|Times|elapsed time|
|---|---|---|---|
|BoltDB|128K|2000|6133872ns|
|MDBM(Store)|256M|3000000|5377ns|




## Links

- [Yahoo! MDBM](https://github.com/yahoo/mdbm)
- [MDBM::Concept](http://yahoo.github.io/mdbm/guide/concepts.html)
- [MDBM::Build](https://github.com/yahoo/mdbm/blob/master/README.build)
- [MDBM::Document](http://yahoo.github.io/mdbm/)
- [MDBM::FAQ](http://yahoo.github.io/mdbm/guide/faq.html)
- [DBM](https://en.wikipedia.org/wiki/Dbm)
- [BoltDB](https://github.com/boltdb/bolt)
- [PHP-MDBM](https://github.com/torden/php-mdbm)
- [Py-MDBM](https://github.com/torden/py-mdbm)

---

*Please feel free. I hope it is helpful for you*
