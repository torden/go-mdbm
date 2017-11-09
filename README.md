# Go-mdbm

---

*NOT READY FOR USE*


- Go-mdbm is a Go bind to Yahoo! MDBM C API.
- MDBM is a super-fast memory-mapped key/value store.

[![Build Status](https://travis-ci.org/torden/go-strutil.svg?branch=master)](https://travis-ci.org/torden/go-strutil)
[![Go Report Card](https://goreportcard.com/badge/github.com/torden/go-strutil)](https://goreportcard.com/report/github.com/torden/go-strutil)
[![GoDoc](https://godoc.org/github.com/torden/go-strutil?status.svg)](https://godoc.org/github.com/torden/go-strutil)
[![Coverage Status](https://coveralls.io/repos/github/torden/go-strutil/badge.svg?branch=master)](https://coveralls.io/github/torden/go-strutil?branch=master)
[![Go Walker](http://gowalker.org/api/v1/badge)](https://gowalker.org/github.com/torden/go-strutil)
[![GitHub version](https://badge.fury.io/gh/torden%2Fgo-strutil.svg)](https://badge.fury.io/gh/torden%2Fgo-strutil)

```
On Now, Almost MDBM APIs supported
```


## Install

### MDBM

#### Download

Use the master branch
```
git clone https://github.com/yahoo/mdbm.git
```
OR Use the release tarball
```
wget https://github.com/yahoo/mdbm/archive/v4.12.3.tar.gz
tar xvzf v4.12.3.tar.gz
```

#### Install

Refer to the https://github.com/yahoo/mdbm/blob/master/README.build

if you want to install to another path, (HIGH RECOMMEND)
```
cd mdbm
PREFIX=/usr/local/mdbm make install
```


### go-mdbm

#### Download for the MDBM installed another path

```
git get github.com/torden/go-mdbm.git
```

Maybe you know, the mdbm installation default path is /tmp/install, That's why go-mdbm used it.
if you did install to another path, You must  change the mdbm installed path in *mdbm.go* source code

```shell
cd $GOPATH/torden/go-mdbm/
vi mdbm.go
```

```go
#cgo CFLAGS: -I/tmp/install/include/
#cgo LDFLAGS: -L/tmp/install/lib64/ -Wl,-rpath=/tmp/install/lib64/ -lmdbm
```

#### Download for the MDBM installed default path

```
git get github.com/yahoo/mdbm.git
```


## Support two compatibility branches:

|*Branch*|*Support*|*test*|
|---|---|---|
|master|yes|always automatic testing|
|release 4.3.x|yes|tested|



## On Now, Not Support APIs

|*API*|*STATUS*|
|---|---|
|mdbm_save|DEPRECATED|
|mdbm_restore|DEPRECATED|
|mdbm_sethash|DEPRECATED|
|mdbm_stat_all_page|as soon|
|mdbm_stat_header|as soon|
|mdbm_cdbdump_add_record|as soon|
|mdbm_cdbdump_import|as soon|
|mdbm_dbdump_to_file|as soon|
|mdbm_dbdump_trailer_and_close|as soon|
|mdbm_dbdump_export_header|as soon|
|mdbm_set_backingstore|as soon|
|mdbm_replace_backing_store|as soon|
|mdbm_pre_split|as soon|
|mdbm_fcopy|as soon|
|mdbm_iterate|as soon|
|mdbm_prune|as soon|
|mdbm_set_cleanfunc|as soon|
|mdbm_clean|as soon|
|mdbm_set_stats_func|as soon|
|mdbm_chunk_iterate|as soon|
|mdbm_sparsify_file|as soon|


## Todo

* Write UnitTests
* Binding All APIs
* Testing on another platform (osx, bsd...)
* EasyXXX APIs

---
Please feel free.
