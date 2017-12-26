# Go-mdbm

*Unfortunately, Not ready for use on production, but go-mdbm tests complete as soon...*

- Go-mdbm is a Go(golang) bind to Yahoo! MDBM C API.
- MDBM is a super-fast memory-mapped key/value store.
- To use the modern buzzwords, it is NoSQL, and for many operations, it is Zero-Copy.
- It is based on an earlier version by Larry McVoy, then at SGI. That in turn, is based on SDBM by Ozan Yigit. wikipedia DBM article
- Yahoo added significant performance enhancements, many tools, tests, and comprehensive documentation.
- It has been used in production for over a decade, for a wide variety of applications, both large and small.

[![Build Status](https://travis-ci.org/torden/go-mdbm.svg?branch=master)](https://travis-ci.org/torden/go-mdbm)
[![Go Report Card](https://goreportcard.com/badge/github.com/torden/go-mdbm)](https://goreportcard.com/report/github.com/torden/go-mdbm)
[![GoDoc](https://godoc.org/github.com/torden/go-mdbm?status.svg)](https://godoc.org/github.com/torden/go-mdbm)
[![Coverage Status](https://coveralls.io/repos/github/torden/go-mdbm/badge.svg?branch=master)](https://coveralls.io/github/torden/go-mdbm?branch=master)
[![Go Walker](http://gowalker.org/api/v1/badge)](https://gowalker.org/github.com/torden/go-mdbm)
[![GitHub version](https://badge.fury.io/gh/torden%2Fgo-mdbm.svg)](https://badge.fury.io/gh/torden%2Fgo-mdbm)

```
Allmost MDBM APIs supported
```

## Install from Source Code

### MDBM

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

```
git get github.com/torden/go-mdbm.git
```

As you know, The mdbm installation default path is /tmp/install, That's why go-mdbm used it.
if you did mdbm install to another path or saw the following message, You must change the mdbm installed path in *mdbm.go* source code

```shell
$ go get github.com/torden/go-mdbm
# github.com/torden/go-mdbm
go-mdbm/mdbm.go:13:10: fatal error: mdbm.h: No such file or directory
 #include <mdbm.h>
          ^~~~~~~~
compilation terminated.
```

#### Change the mdbm installed path 

if you did change any other installation path, you must following below

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

## On Now, Not Support APIs


### Deprecatred APIs

|*API*|*STATUS*|*COMMENT*|
|---|---|---|
|mdbm_save|DEPRECATED|mdbm_save is only supported for V2 MDBMs.|
|mdbm_restore|DEPRECATED|mdbm_restore is only supported for V2 MDBMs.|
|mdbm_sethash|DEPRECATED|Legacy version of mdbm_set_hash() This function has inconsistent naming, and error return value. It will be removed in a future version.|

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

# Example

[Source Code and a sample file](https://github.com/torden/go-mdbm/tree/master/example)

	
## Todo

* coverage up to 95% (min.)
* Binding All APIs without deprecated apis
* Testing on another platform (osx, bsd...)
* Pre-compiled mdbm library by OS
* Stabilization & Clear

---
Please feel free.
