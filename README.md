# Go-mdbm

- Go-mdbm is a Go(golang,go-lang) bind to [Yahoo! MDBM C API](https://github.com/yahoo/mdbm).
- MDBM is a super-fast memory-mapped key/value store.

[![Build Status](https://travis-ci.org/torden/go-mdbm.svg?branch=master)](https://travis-ci.org/torden/go-mdbm)
[![Go Report Card](https://goreportcard.com/badge/github.com/torden/go-mdbm)](https://goreportcard.com/report/github.com/torden/go-mdbm)
[![GoDoc](https://godoc.org/github.com/torden/go-mdbm?status.svg)](https://godoc.org/github.com/torden/go-mdbm)
[![Coverage Status](https://coveralls.io/repos/github/torden/go-mdbm/badge.svg?branch=master)](https://coveralls.io/github/torden/go-mdbm?branch=master)
[![Go Walker](http://gowalker.org/api/v1/badge)](https://gowalker.org/github.com/torden/go-mdbm)
[![GitHub version](https://badge.fury.io/gh/torden%2Fgo-mdbm.svg)](https://badge.fury.io/gh/torden%2Fgo-mdbm)

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

#### Download

for use

```
git get -u github.com/torden/go-mdbm
```

for development or customization

```
git clone https://github.com/torden/go-mdbm
```

#### Preparing to Build

As you know, The mdbm installation default path is /tmp/install, That's why go-mdbm used it.
if you did the mdbm install to another path or saw the following message, You must change the mdbm installed path in *mdbm.go* source code

```shell
$ go get -u github.com/torden/go-mdbm
# github.com/torden/go-mdbm
go-mdbm/mdbm.go:13:10: fatal error: mdbm.h: No such file or directory
 #include <mdbm.h>
          ^~~~~~~~
compilation terminated.
```

#### Build

```shell
cd $GOPATH/src/github.com/torden/go-mdbm

make clean
make setup
make build
```

#### Testing

```shell
cd $GOPATH/src/github.com/torden/go-mdbm
make test
```

#### Run to Example

```shell
cd $GOPATH/src/github.com/torden/go-mdbm/example/
go run -race example.go
```

#### Miscellaneous


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

#### Change the mdbm installed path

if you did change the mdbm installation path, you must be following.

```shell
cd $GOPATH/src/github.com/torden/go-mdbm/
vi mdbm.go
```

```go
#cgo CFLAGS: -I/[MDBM_INSTALLED_PATH]/include/
#cgo LDFLAGS: -L/[MDBM_INSTALLED_PATH]/lib64/ -Wl,-rpath=/[MDBM_INSTALLED_PATH]/lib64/ -lmdbm
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


## Benchmark

The following is result of Go-mdbm vs  BoltDB benchmarks for simple data storing and random fetching in them.

### Command

[Source Code](https://github.com/torden/go-mdbm/blob/master/benchmark_test.go)

```shell
go test -race -bench=. -run Benchmark -test.benchmem -v -cpu 8
```

### Output

|Func|Times|NanoSecond per Loop|Byte|Allocs|
|---|---|---|---|---|
|Benchmark_boltdb_Store-8|200|7917647ns/op|29052B/op|52allocs/op|
|Benchmark_boltdb_Fetch-8|200000|14746ns/op|496B/op|9allocs/op|
|Benchmark_mdbm_Store-8|300000|3764ns/op|96B/op|6allocs/op|
|Benchmark_mdbm_StoreWithLock-8|500000|3768ns/op|96B/op|6allocs/op|
|Benchmark_mdbm_StoreOnLock-8|500000|3473ns/op|96B/op|6allocs/op|
|Benchmark_mdbm_Fetch-8|1000000|2200ns/op|56B/op|4allocs/op|
|Benchmark_mdbm_FetchWithLock-8|500000|2751ns/op|56B/op|4allocs/op|
|Benchmark_mdbm_FetchOnLock-8|500000|2785ns/op|56B/op|4allocs/op|
|Benchmark_mdbm_PreLoad_Fetch-8|500000|2150ns/op|56B/op|4allocs/op|
|Benchmark_mdbm_PreLoad_FetchWithLock-8|500000|2762ns/op|56B/op|4allocs/op|
|Benchmark_mdbm_PreLoad_FetchOnLock-8|500000|2620ns/op|56B/op|4allocs/op|



## Links

- [Yahoo! MDBM](https://github.com/yahoo/mdbm)
- [MDBM Build](https://github.com/yahoo/mdbm/blob/master/README.build)
- [MDBM Document](http://yahoo.github.io/mdbm/)
- [DBM](https://en.wikipedia.org/wiki/Dbm)
- [BoltDB](https://github.com/boltdb/bolt)

---
Please feel free. I hope it is helpful for you
