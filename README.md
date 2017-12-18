# Go-mdbm

*Not ready for use on production servers, go-mdbm tests complete as soon...*

- Go-mdbm is a Go bind to Yahoo! MDBM C API.
- MDBM is a super-fast memory-mapped key/value store.

[![Build Status](https://travis-ci.org/torden/go-mdbm.svg?branch=master)](https://travis-ci.org/torden/go-mdbm)
[![Go Report Card](https://goreportcard.com/badge/github.com/torden/go-mdbm)](https://goreportcard.com/report/github.com/torden/go-mdbm)
[![GoDoc](https://godoc.org/github.com/torden/go-mdbm?status.svg)](https://godoc.org/github.com/torden/go-mdbm)
[![Coverage Status](https://coveralls.io/repos/github/torden/go-mdbm/badge.svg?branch=master)](https://coveralls.io/github/torden/go-mdbm?branch=master)
[![Go Walker](http://gowalker.org/api/v1/badge)](https://gowalker.org/github.com/torden/go-mdbm)
[![GitHub version](https://badge.fury.io/gh/torden%2Fgo-mdbm.svg)](https://badge.fury.io/gh/torden%2Fgo-mdbm)

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

```shell
cd mdbm
PREFIX=/usr/local/mdbm make install
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
|master|yes|always testing|
|release 4.3.x|yes|always testing|

## Support two compatibility O.S 

|*Dist*|*Support*|*test*|
|---|---|---|
|rhel|yes|always testing|
|ubuntu|yes|always testing|
|bsd|as soon|as soon|
|osx|as soon|as soon|

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
|mdbm_fetch_buf|similar Fetch() api|

# Example

[Full Source Code and a sample file](https://github.com/torden/go-mdbm/tree/master/example)


## Creating and populating a database

### Code
```go
package main

import (
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	mdbm "github.com/torden/go-mdbm"
)

const (
	mdbmPath1   = "/tmp/example1.mdbm"
	mdbmPath2   = "/tmp/example2.mdbm"
	sample1Path = "./sample1.tsv" //ISO Alpha-2,3 and Numeric Country Codes
)

func main() {

	var err error
	err = os.Remove(mdbmPath1)
	if err == nil {
		log.Println("not exists the mdbm file(=%s)", mdbmPath1)
	} else {
		log.Println("remove the mdbm file(=%s)", mdbmPath1)
	}
    
	log.Printf("Createing and Populating the mdbm file(=%s)", mdbmPath1)

	//init. the go-mdbm
	dbm := mdbm.NewMDBM()

	//create & open(RDRW) an mdbm file
	err = dbm.EasyOpen(mdbmPath1, 0644)

	//the mdbm object close at close func
	defer dbm.EasyClose()

	//check the open error
	if err != nil {
		log.Fatalf("failed, can't open mdbm file\npath=%s, err=%v", mdbmPath1, err)
	}

	//read a content
	data, err := ioutil.ReadFile(sample1Path)
	if err != nil {
		log.Fatalf("failed, read to the tsv file\npath=%s, err=%v", sample1Path, err)
	}

	//convert []byte to string and split by newline
	dataStrAr := strings.Split(string(data), "\n")
	for k, v := range dataStrAr {

		if k == 0 { //header
			continue
		}

		//obtain data by field
		row := strings.Split(v, "\t")
		if len(row) < 4 {
			continue
		}

		//0 : Country or Area Name
		//1 : ISO ALPHA-2 Code
		//2 : ISO ALPHA-3 Code
		//3 : USO Numeric Code , UN M49 Numeric Code
		rv, err := dbm.Store(row[0], row[1], mdbm.Insert) // if key does not exist; fail if exists
		if err != nil {
			log.Fatalf("failed, can't data(key=%+v, value=%+v) add to the mdbm file(=%s)\nrv=%d, err=%v", row[0], row[1], mdbmPath1, rv, err)
		}
	}

	log.Println("complete")
}

```

### Output
```
2017/01/01 00:00:00 remove the mdbm file(=%s) /tmp/example1.mdbm
2017/01/01 00:00:00 not exists the mdbm file(=%s) /tmp/example2.mdbm
2017/01/01 00:00:00 Createing and Populating the mdbm file(=/tmp/example1.mdbm)
2017/01/01 00:00:00 complete
```


## Creating and populating a database

### Code
```go
package main

import (
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	mdbm "github.com/torden/go-mdbm"
)

const (
	mdbmPath1   = "/tmp/example1.mdbm"
	mdbmPath2   = "/tmp/example2.mdbm"
	sample1Path = "./sample1.tsv" //ISO Alpha-2,3 and Numeric Country Codes
)

func main() {

	var rv int
	var err error

	err = os.Remove(mdbmPath2)
		if err == nil {
		log.Println("not exists the mdbm file(=%s)", mdbmPath2)
	} else {
		log.Println("remove the mdbm file(=%s)", mdbmPath2)
	}

	log.Printf("Createing and Populating the mdbm file(=%s)", mdbmPath2)

	//init. the go-mdbm
	dbm := mdbm.NewMDBM()

	//create & open(RDRW) an mdbm file
	err = dbm.EasyOpen(mdbmPath2, 0644)

	//the mdbm object close at close func
	defer dbm.EasyClose()

	//check the open error
	if err != nil {
		log.Fatalf("failed, can't open mdbm file\npath=%s, err=%v", mdbmPath2, err)
	}

	vbyte := byte(60)
	varbyte := []byte("hello")
	vstring := "hello_world"

	vint := int(1)
	//vint8 := int8(12) //not support. because int8 and byte are same data type
	vint16 := int16(3)
	vint32 := int32(4)
	vint64 := int64(5)

	vuint := uint(11)
	//vuint8 := uint8(12) //not support. because uint8 and byte are same data type
	vuint16 := uint16(13)
	vuint32 := uint32(14)
	vuint64 := uint64(15)

	vfloat32 := float32(2.1)
	vfloat64 := float64(2.2)

	vcomplex64 := complex64(3)
	vcomplex128 := complex128(4)

	//bool
	//byte
	//[]uint8
	//string
	rv, err = dbm.StoreWithLock(vbyte, vbyte, mdbm.Insert) // if key does not exist; fail if exists
	if err != nil {
		log.Fatalf("failed, can't data(byte) add to the mdbm file(=%s)\nrv=%d, err=%v", mdbmPath2, rv, err)
	}

	rv, err = dbm.StoreWithLock(varbyte, varbyte, mdbm.Insert) // if key does not exist; fail if exists
	if err != nil {
		log.Fatalf("failed, can't data([]byte) add to the mdbm file(=%s)\nrv=%d, err=%v", mdbmPath2, rv, err)
	}

	rv, err = dbm.StoreWithLock(vstring, vstring, mdbm.Insert) // if key does not exist; fail if exists
	if err != nil {
		log.Fatalf("failed, can't data(string) add to the mdbm file(=%s)\nrv=%d, err=%v", mdbmPath2, rv, err)
	}

	rv, err = dbm.StoreWithLock(vint, vint, mdbm.Insert) // if key does not exist; fail if exists
	if err != nil {
		log.Fatalf("failed, can't data(int) add to the mdbm file(=%s)\nrv=%d, err=%v", mdbmPath2, rv, err)
	}

	rv, err = dbm.StoreWithLock(vint16, vint16, mdbm.Insert) // if key does not exist; fail if exists
	if err != nil {
		log.Fatalf("failed, can't data(int16) add to the mdbm file(=%s)\nrv=%d, err=%v", mdbmPath2, rv, err)
	}

	rv, err = dbm.StoreWithLock(vint32, vint32, mdbm.Insert) // if key does not exist; fail if exists
	if err != nil {
		log.Fatalf("failed, can't data(int32) add to the mdbm file(=%s)\nrv=%d, err=%v", mdbmPath2, rv, err)
	}

	rv, err = dbm.StoreWithLock(vint64, vint64, mdbm.Insert) // if key does not exist; fail if exists
	if err != nil {
		log.Fatalf("failed, can't data(int32) add to the mdbm file(=%s)\nrv=%d, err=%v", mdbmPath2, rv, err)
	}

	rv, err = dbm.StoreWithLock(vuint, vuint, mdbm.Insert) // if key does not exist; fail if exists
	if err != nil {
		log.Fatalf("failed, can't data(uint) add to the mdbm file(=%s)\nrv=%d, err=%v", mdbmPath2, rv, err)
	}

	rv, err = dbm.StoreWithLock(vuint16, vuint16, mdbm.Insert) // if key does not exist; fail if exists
	if err != nil {
		log.Fatalf("failed, can't data(uint16) add to the mdbm file(=%s)\nrv=%d, err=%v", mdbmPath2, rv, err)
	}

	rv, err = dbm.StoreWithLock(vuint32, vuint32, mdbm.Insert) // if key does not exist; fail if exists
	if err != nil {
		log.Fatalf("failed, can't data(uint32) add to the mdbm file(=%s)\nrv=%d, err=%v", mdbmPath2, rv, err)
	}

	rv, err = dbm.StoreWithLock(vuint64, vuint64, mdbm.Insert) // if key does not exist; fail if exists
	if err != nil {
		log.Fatalf("failed, can't data(uint32) add to the mdbm file(=%s)\nrv=%d, err=%v", mdbmPath2, rv, err)
	}

	rv, err = dbm.StoreWithLock(vfloat32, vfloat32, mdbm.Insert) // if key does not exist; fail if exists
	if err != nil {
		log.Fatalf("failed, can't data(float32) add to the mdbm file(=%s)\nrv=%d, err=%v", mdbmPath2, rv, err)
	}

	rv, err = dbm.StoreWithLock(vfloat64, vfloat64, mdbm.Insert) // if key does not exist; fail if exists
	if err != nil {
		log.Fatalf("failed, can't data(float64) add to the mdbm file(=%s)\nrv=%d, err=%v", mdbmPath2, rv, err)
	}

	rv, err = dbm.StoreWithLock(vcomplex64, vcomplex64, mdbm.Insert) // if key does not exist; fail if exists
	if err != nil {
		log.Fatalf("failed, can't data(complex64) add to the mdbm file(=%s)\nrv=%d, err=%v", mdbmPath2, rv, err)
	}

	rv, err = dbm.StoreWithLock(vcomplex128, vcomplex128, mdbm.Insert) // if key does not exist; fail if exists
	if err != nil {
		log.Fatalf("failed, can't data(complex128) add to the mdbm file(=%s)\nrv=%d, err=%v", mdbmPath2, rv, err)
	}

	dbm.Sync()

	//validation
	var cnt int
	key, val, err := dbm.First()
	if err != nil {
		log.Fatalf("key=%s, val=%s, err=%v", key, val, err)
	} else {
		log.Printf("key = %-20s | value = %s", key, val)
		cnt++
	}

	for {

		key, val, err := dbm.Next()
		if err != nil {
			log.Fatalf("key=%s, val=%s, err=%v", key, val, err)
		}

		if len(key) < 1 {
			break
		}

		if key != val {
			log.Fatalf("key and valiue mismatch.\nKey=%s , Value=%s", key, val)
		}

		log.Printf("key = %-20s | value = %s", key, val)
		cnt++
	}

	log.Printf("the count of number of rows in the mdbm(=%s) is `%d` rows", mdbmPath2, cnt)
	log.Println("complete")
}
```

### Output
```
2017/01/01 00:00:00 remove the mdbm file(=%s) /tmp/example2.mdbm
2017/01/01 00:00:00 Createing and Populating the mdbm file(=/tmp/example2.mdbm)
2017/01/01 00:00:00 key = <                    | value = <
2017/01/01 00:00:00 key = hello                | value = hello
2017/01/01 00:00:00 key = hello_world          | value = hello_world
2017/01/01 00:00:00 key = 1                    | value = 1
2017/01/01 00:00:00 key = 3                    | value = 3
2017/01/01 00:00:00 key = 4                    | value = 4
2017/01/01 00:00:00 key = 5                    | value = 5
2017/01/01 00:00:00 key = 11                   | value = 11
2017/01/01 00:00:00 key = 13                   | value = 13
2017/01/01 00:00:00 key = 14                   | value = 14
2017/01/01 00:00:00 key = 15                   | value = 15
2017/01/01 00:00:00 key = 2.1                  | value = 2.1
2017/01/01 00:00:00 key = 2.2                  | value = 2.2
2017/01/01 00:00:00 key = (3+0i)               | value = (3+0i)
2017/01/01 00:00:00 key = (4+0i)               | value = (4+0i)
2017/01/01 00:00:00 the count of number of rows in the mdbm(=/tmp/example2.mdbm) is `15` rows
2017/01/01 00:00:00 complete
```


## Fetching records in-place

### Code

```go
package main

import (
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	mdbm "github.com/torden/go-mdbm"
)

const (
	mdbmPath1   = "/tmp/example1.mdbm"
	mdbmPath2   = "/tmp/example2.mdbm"
	sample1Path = "./sample1.tsv" //ISO Alpha-2,3 and Numeric Country Codes
)

func main() {

	var err error

	log.Printf("Fetcing records in the mdbm file(=%s)", mdbmPath1)

	//init. the go-mdbm
	dbm := mdbm.NewMDBM()

	//create & open(RDRW) an mdbm file
	err = dbm.EasyOpen(mdbmPath1, 0644)

	//the mdbm object close at close func
	defer dbm.EasyClose()

	//check the open error
	if err != nil {
		log.Fatalf("failed, can't open mdbm file\npath=%s, err=%v", mdbmPath1, err)
	}

	//read a content
	data, err := ioutil.ReadFile(sample1Path)
	if err != nil {
		log.Fatalf("failed, read to the tsv file\npath=%s, err=%v", sample1Path, err)
	}

	//convert []byte to string and split by newline
	dataStrAr := strings.Split(string(data), "\n")

	var keyList []string
	var keySize int

	//obtain list of key
	for k, v := range dataStrAr {

		if k == 0 { //header
			continue
		}

		row := strings.Split(v, "\t")
		if len(row) < 4 {
			continue
		}

		keyList = append(keyList, row[0])
	}

	keySize = len(keyList)

	//obtain a pseudo-random 32-bit value as a uint32
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	var rkey int

	//random fetching..
	for i := 0; i <= 100; i++ {

		//obtain a pseudo-random 32-bit value between 0 and keySize
		for {
			rkey = random.Intn(keySize)
			if rkey >= 0 {
				break
			}
		}

		//fetch on locking
		val, err := dbm.FetchWithLock(keyList[rkey])
		if err != nil {
			log.Fatalf("failed, can't find out value in the mdbm file(=%s)\nekey=%s, err=%v", mdbmPath1, keyList[rkey], err)
		}

		log.Printf("ISO ALPHA-2 Code : [%s] , Country or Area Name : [%s]", keyList[rkey], val)
	}

	log.Println("complete")
}
```

### Output
```
2017/01/01 00:00:00 Fetcing records in the mdbm file(=/tmp/example1.mdbm)
2017/01/01 00:00:00 ISO ALPHA-2 Code : [TL] , Country or Area Name : [Timor-Leste]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [KH] , Country or Area Name : [Cambodia]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [TR] , Country or Area Name : [Turkey]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [OM] , Country or Area Name : [Oman]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [BS] , Country or Area Name : [Bahamas]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [MV] , Country or Area Name : [Maldives]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [BR] , Country or Area Name : [Brazil]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [PY] , Country or Area Name : [Paraguay]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [JO] , Country or Area Name : [Jordan]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [KP] , Country or Area Name : [Korea (North)]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [MG] , Country or Area Name : [Madagascar]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [PA] , Country or Area Name : [Panama]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [BR] , Country or Area Name : [Brazil]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [MT] , Country or Area Name : [Malta]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [AX] , Country or Area Name : [Aland Islands]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [TL] , Country or Area Name : [Timor-Leste]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [PG] , Country or Area Name : [Papua New Guinea]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [CY] , Country or Area Name : [Cyprus]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [MD] , Country or Area Name : [Moldova]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [GU] , Country or Area Name : [Guam]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [CI] , Country or Area Name : [Côte d'Ivoire]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [KP] , Country or Area Name : [Korea (North)]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [ZA] , Country or Area Name : [South Africa]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [LI] , Country or Area Name : [Liechtenstein]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [AE] , Country or Area Name : [United Arab Emirates]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [NU] , Country or Area Name : [Niue]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [TR] , Country or Area Name : [Turkey]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [BJ] , Country or Area Name : [Benin]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [AQ] , Country or Area Name : [Antarctica]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [AR] , Country or Area Name : [Argentina]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [DO] , Country or Area Name : [Dominican Republic]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [AI] , Country or Area Name : [Anguilla]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [GR] , Country or Area Name : [Greece]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [AE] , Country or Area Name : [United Arab Emirates]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [IM] , Country or Area Name : [Isle of Man]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [KE] , Country or Area Name : [Kenya]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [VA] , Country or Area Name : [Holy See (Vatican City State)]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [TC] , Country or Area Name : [Turks and Caicos Islands]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [WS] , Country or Area Name : [Samoa]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [KM] , Country or Area Name : [Comoros]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [AT] , Country or Area Name : [Austria]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [DK] , Country or Area Name : [Denmark]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [DM] , Country or Area Name : [Dominica]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [VI] , Country or Area Name : [Virgin Islands, US]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [NC] , Country or Area Name : [New Caledonia]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [SH] , Country or Area Name : [Saint Helena]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [GH] , Country or Area Name : [Ghana]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [MH] , Country or Area Name : [Marshall Islands]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [GL] , Country or Area Name : [Greenland]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [PS] , Country or Area Name : [Palestinian Territory]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [CC] , Country or Area Name : [Cocos (Keeling) Islands]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [BD] , Country or Area Name : [Bangladesh]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [EE] , Country or Area Name : [Estonia]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [IS] , Country or Area Name : [Iceland]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [CO] , Country or Area Name : [Colombia]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [GE] , Country or Area Name : [Georgia]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [BI] , Country or Area Name : [Burundi]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [ZA] , Country or Area Name : [South Africa]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [DM] , Country or Area Name : [Dominica]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [SR] , Country or Area Name : [Suriname]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [PL] , Country or Area Name : [Poland]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [TF] , Country or Area Name : [French Southern Territories]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [SC] , Country or Area Name : [Seychelles]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [AM] , Country or Area Name : [Armenia]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [CY] , Country or Area Name : [Cyprus]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [NL] , Country or Area Name : [Netherlands]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [IE] , Country or Area Name : [Ireland]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [LK] , Country or Area Name : [Sri Lanka]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [JM] , Country or Area Name : [Jamaica]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [SC] , Country or Area Name : [Seychelles]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [CZ] , Country or Area Name : [Czech Republic]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [FK] , Country or Area Name : [Falkland Islands (Malvinas)]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [JO] , Country or Area Name : [Jordan]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [MN] , Country or Area Name : [Mongolia]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [TF] , Country or Area Name : [French Southern Territories]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [TW] , Country or Area Name : [Taiwan, Republic of China]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [AF] , Country or Area Name : [Afghanistan]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [UY] , Country or Area Name : [Uruguay]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [AO] , Country or Area Name : [Angola]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [VG] , Country or Area Name : [British Virgin Islands]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [FI] , Country or Area Name : [Finland]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [RO] , Country or Area Name : [Romania]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [LA] , Country or Area Name : [Lao PDR]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [GG] , Country or Area Name : [Guernsey]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [MF] , Country or Area Name : [Saint-Martin (French part)]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [LV] , Country or Area Name : [Latvia]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [AZ] , Country or Area Name : [Azerbaijan]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [YE] , Country or Area Name : [Yemen]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [FJ] , Country or Area Name : [Fiji]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [PE] , Country or Area Name : [Peru]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [PR] , Country or Area Name : [Puerto Rico]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [BO] , Country or Area Name : [Bolivia]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [VN] , Country or Area Name : [Viet Nam]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [MH] , Country or Area Name : [Marshall Islands]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [BO] , Country or Area Name : [Bolivia]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [TJ] , Country or Area Name : [Tajikistan]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [PF] , Country or Area Name : [French Polynesia]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [GA] , Country or Area Name : [Gabon]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [FK] , Country or Area Name : [Falkland Islands (Malvinas)]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [KW] , Country or Area Name : [Kuwait]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [ZW] , Country or Area Name : [Zimbabwe]
2017/01/01 00:00:00 complete
```



## Iterating over all records¶

### Code
```go
package main

import (
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	mdbm "github.com/torden/go-mdbm"
)

const (
	mdbmPath1   = "/tmp/example1.mdbm"
	mdbmPath2   = "/tmp/example2.mdbm"
	sample1Path = "./sample1.tsv" //ISO Alpha-2,3 and Numeric Country Codes
)

func main() {

	var err error

	log.Printf("Iterating over all records in the mdbm file(=%s)", mdbmPath1)

	//init. the go-mdbm
	dbm := mdbm.NewMDBM()

	//create & open(RDRW) an mdbm file
	err = dbm.EasyOpen(mdbmPath1, 0644)

	//the mdbm object close at close func
	defer dbm.EasyClose()

	//check the open error
	if err != nil {
		log.Fatalf("failed, can't open mdbm file\npath=%s, err=%v", mdbmPath1, err)
	}

	key, val, err := dbm.First()
	if err != nil {
		log.Fatalf("key=%s, val=%s, err=%v", key, val, err)
	}

	log.Printf("ISO ALPHA-2 Code : [%s] , Country or Area Name : [%s]", keyList[rkey], val)

	var cnt uint64

	for {

		key, val, err := dbm.Next()
		if len(key) < 1 {
			break
		}

		if err != nil {
			log.Fatalf("key=%s, val=%s, err=%v", key, val, err)
		}

		log.Printf("ISO ALPHA-2 Code : [%s] , Country or Area Name : [%s]", keyList[rkey], val)
		cnt++
	}

	log.Printf("the count of number of rows in the mdbm(=%s) is `%d` rows", mdbmPath1, cnt)
	log.Println("complete")
}

```

### Output
```
2017/01/01 00:00:00 Iterating over all records in the mdbm file(=/tmp/example1.mdbm)
2017/01/01 00:00:00 ISO ALPHA-2 Code : [AF] , Country or Area Name : [Afghanistan]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [AX] , Country or Area Name : [Aland Islands]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [AL] , Country or Area Name : [Albania]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [AD] , Country or Area Name : [Andorra]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [AR] , Country or Area Name : [Argentina]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [AT] , Country or Area Name : [Austria]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [AZ] , Country or Area Name : [Azerbaijan]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [BS] , Country or Area Name : [Bahamas]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [BY] , Country or Area Name : [Belarus]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [BE] , Country or Area Name : [Belgium]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [BM] , Country or Area Name : [Bermuda]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [BO] , Country or Area Name : [Bolivia]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [BA] , Country or Area Name : [Bosnia and Herzegovina]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [BW] , Country or Area Name : [Botswana]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [VG] , Country or Area Name : [British Virgin Islands]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [BG] , Country or Area Name : [Bulgaria]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [ZA] , Country or Area Name : [South Africa]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [BI] , Country or Area Name : [Burundi]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [KH] , Country or Area Name : [Cambodia]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [CV] , Country or Area Name : [Cape Verde]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [CF] , Country or Area Name : [Central African Republic]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [CL] , Country or Area Name : [Chile]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [CN] , Country or Area Name : [China]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [HK] , Country or Area Name : [Hong Kong, SAR China]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [CX] , Country or Area Name : [Christmas Island]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [CD] , Country or Area Name : [Congo, (Kinshasa)]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [CR] , Country or Area Name : [Costa Rica]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [CZ] , Country or Area Name : [Czech Republic]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [DK] , Country or Area Name : [Denmark]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [SR] , Country or Area Name : [Suriname]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [DM] , Country or Area Name : [Dominica]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [DO] , Country or Area Name : [Dominican Republic]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [SV] , Country or Area Name : [El Salvador]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [ER] , Country or Area Name : [Eritrea]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [ET] , Country or Area Name : [Ethiopia]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [FK] , Country or Area Name : [Falkland Islands (Malvinas)]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [FO] , Country or Area Name : [Faroe Islands]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [FI] , Country or Area Name : [Finland]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [GF] , Country or Area Name : [French Guiana]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [DE] , Country or Area Name : [Germany]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [GH] , Country or Area Name : [Ghana]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [GR] , Country or Area Name : [Greece]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [GL] , Country or Area Name : [Greenland]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [GD] , Country or Area Name : [Grenada]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [GP] , Country or Area Name : [Guadeloupe]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [GT] , Country or Area Name : [Guatemala]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [GN] , Country or Area Name : [Guinea]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [HM] , Country or Area Name : [Heard and Mcdonald Islands]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [VA] , Country or Area Name : [Holy See (Vatican City State)]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [HU] , Country or Area Name : [Hungary]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [IN] , Country or Area Name : [India]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [ID] , Country or Area Name : [Indonesia]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [IR] , Country or Area Name : [Iran, Islamic Republic of]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [IL] , Country or Area Name : [Israel]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [IT] , Country or Area Name : [Italy]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [JM] , Country or Area Name : [Jamaica]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [SD] , Country or Area Name : [Sudan]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [JE] , Country or Area Name : [Jersey]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [JO] , Country or Area Name : [Jordan]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [KZ] , Country or Area Name : [Kazakhstan]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [KP] , Country or Area Name : [Korea (North)]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [KR] , Country or Area Name : [Korea (South)]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [LA] , Country or Area Name : [Lao PDR]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [LS] , Country or Area Name : [Lesotho]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [LY] , Country or Area Name : [Libya]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [LI] , Country or Area Name : [Liechtenstein]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [LU] , Country or Area Name : [Luxembourg]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [MV] , Country or Area Name : [Maldives]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [ML] , Country or Area Name : [Mali]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [MT] , Country or Area Name : [Malta]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [MH] , Country or Area Name : [Marshall Islands]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [MR] , Country or Area Name : [Mauritania]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [YT] , Country or Area Name : [Mayotte]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [MX] , Country or Area Name : [Mexico]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [FM] , Country or Area Name : [Micronesia, Federated States of]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [MD] , Country or Area Name : [Moldova]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [MN] , Country or Area Name : [Mongolia]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [MZ] , Country or Area Name : [Mozambique]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [NA] , Country or Area Name : [Namibia]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [AN] , Country or Area Name : [Netherlands Antilles]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [NC] , Country or Area Name : [New Caledonia]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [NI] , Country or Area Name : [Nicaragua]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [NE] , Country or Area Name : [Niger]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [NG] , Country or Area Name : [Nigeria]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [NU] , Country or Area Name : [Niue]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [MP] , Country or Area Name : [Northern Mariana Islands]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [NO] , Country or Area Name : [Norway]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [PK] , Country or Area Name : [Pakistan]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [PW] , Country or Area Name : [Palau]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [PS] , Country or Area Name : [Palestinian Territory]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [PA] , Country or Area Name : [Panama]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [PG] , Country or Area Name : [Papua New Guinea]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [PY] , Country or Area Name : [Paraguay]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [PE] , Country or Area Name : [Peru]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [RE] , Country or Area Name : [Réunion]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [RO] , Country or Area Name : [Romania]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [RU] , Country or Area Name : [Russian Federation]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [RW] , Country or Area Name : [Rwanda]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [SH] , Country or Area Name : [Saint Helena]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [KN] , Country or Area Name : [Saint Kitts and Nevis]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [LC] , Country or Area Name : [Saint Lucia]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [MF] , Country or Area Name : [Saint-Martin (French part)]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [PM] , Country or Area Name : [Saint Pierre and Miquelon]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [VC] , Country or Area Name : [Saint Vincent and Grenadines]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [ST] , Country or Area Name : [Sao Tome and Principe]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [SN] , Country or Area Name : [Senegal]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [RS] , Country or Area Name : [Serbia]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [SL] , Country or Area Name : [Sierra Leone]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [SB] , Country or Area Name : [Solomon Islands]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [LK] , Country or Area Name : [Sri Lanka]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [SJ] , Country or Area Name : [Svalbard and Jan Mayen Islands]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [SZ] , Country or Area Name : [Swaziland]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [CH] , Country or Area Name : [Switzerland]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [TW] , Country or Area Name : [Taiwan, Republic of China]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [TG] , Country or Area Name : [Togo]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [TK] , Country or Area Name : [Tokelau]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [TO] , Country or Area Name : [Tonga]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [TM] , Country or Area Name : [Turkmenistan]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [TC] , Country or Area Name : [Turks and Caicos Islands]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [GB] , Country or Area Name : [United Kingdom]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [UZ] , Country or Area Name : [Uzbekistan]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [VU] , Country or Area Name : [Vanuatu]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [VE] , Country or Area Name : [Venezuela (Bolivarian Republic)]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [VI] , Country or Area Name : [Virgin Islands, US]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [WF] , Country or Area Name : [Wallis and Futuna Islands]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [EH] , Country or Area Name : [Western Sahara]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [ZM] , Country or Area Name : [Zambia]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [ZW] , Country or Area Name : [Zimbabwe]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [DZ] , Country or Area Name : [Algeria]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [AS] , Country or Area Name : [American Samoa]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [AO] , Country or Area Name : [Angola]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [AI] , Country or Area Name : [Anguilla]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [AQ] , Country or Area Name : [Antarctica]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [AG] , Country or Area Name : [Antigua and Barbuda]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [AM] , Country or Area Name : [Armenia]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [AW] , Country or Area Name : [Aruba]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [AU] , Country or Area Name : [Australia]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [BH] , Country or Area Name : [Bahrain]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [BD] , Country or Area Name : [Bangladesh]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [BB] , Country or Area Name : [Barbados]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [BZ] , Country or Area Name : [Belize]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [BJ] , Country or Area Name : [Benin]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [BT] , Country or Area Name : [Bhutan]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [BV] , Country or Area Name : [Bouvet Island]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [BR] , Country or Area Name : [Brazil]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [IO] , Country or Area Name : [British Indian Ocean Territory]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [BN] , Country or Area Name : [Brunei Darussalam]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [BF] , Country or Area Name : [Burkina Faso]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [CM] , Country or Area Name : [Cameroon]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [CA] , Country or Area Name : [Canada]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [KY] , Country or Area Name : [Cayman Islands]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [TD] , Country or Area Name : [Chad]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [MO] , Country or Area Name : [Macao, SAR China]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [CC] , Country or Area Name : [Cocos (Keeling) Islands]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [CO] , Country or Area Name : [Colombia]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [KM] , Country or Area Name : [Comoros]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [CG] , Country or Area Name : [Congo (Brazzaville)]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [CK] , Country or Area Name : [Cook Islands]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [CI] , Country or Area Name : [Côte d'Ivoire]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [HR] , Country or Area Name : [Croatia]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [CU] , Country or Area Name : [Cuba]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [CY] , Country or Area Name : [Cyprus]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [DJ] , Country or Area Name : [Djibouti]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [EC] , Country or Area Name : [Ecuador]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [EG] , Country or Area Name : [Egypt]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [GQ] , Country or Area Name : [Equatorial Guinea]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [EE] , Country or Area Name : [Estonia]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [FJ] , Country or Area Name : [Fiji]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [FR] , Country or Area Name : [France]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [PF] , Country or Area Name : [French Polynesia]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [TF] , Country or Area Name : [French Southern Territories]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [GA] , Country or Area Name : [Gabon]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [GM] , Country or Area Name : [Gambia]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [GE] , Country or Area Name : [Georgia]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [GI] , Country or Area Name : [Gibraltar]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [GU] , Country or Area Name : [Guam]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [GG] , Country or Area Name : [Guernsey]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [GW] , Country or Area Name : [Guinea-Bissau]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [GY] , Country or Area Name : [Guyana]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [HT] , Country or Area Name : [Haiti]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [HN] , Country or Area Name : [Honduras]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [IS] , Country or Area Name : [Iceland]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [IQ] , Country or Area Name : [Iraq]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [IE] , Country or Area Name : [Ireland]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [IM] , Country or Area Name : [Isle of Man]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [JP] , Country or Area Name : [Japan]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [KE] , Country or Area Name : [Kenya]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [KI] , Country or Area Name : [Kiribati]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [KW] , Country or Area Name : [Kuwait]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [KG] , Country or Area Name : [Kyrgyzstan]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [LV] , Country or Area Name : [Latvia]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [LB] , Country or Area Name : [Lebanon]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [LR] , Country or Area Name : [Liberia]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [LT] , Country or Area Name : [Lithuania]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [MK] , Country or Area Name : [Macedonia, Republic of]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [MG] , Country or Area Name : [Madagascar]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [MW] , Country or Area Name : [Malawi]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [MY] , Country or Area Name : [Malaysia]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [MQ] , Country or Area Name : [Martinique]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [MU] , Country or Area Name : [Mauritius]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [MC] , Country or Area Name : [Monaco]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [ME] , Country or Area Name : [Montenegro]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [MS] , Country or Area Name : [Montserrat]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [MA] , Country or Area Name : [Morocco]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [MM] , Country or Area Name : [Myanmar]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [NR] , Country or Area Name : [Nauru]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [NP] , Country or Area Name : [Nepal]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [NL] , Country or Area Name : [Netherlands]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [NZ] , Country or Area Name : [New Zealand]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [NF] , Country or Area Name : [Norfolk Island]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [OM] , Country or Area Name : [Oman]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [PH] , Country or Area Name : [Philippines]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [PN] , Country or Area Name : [Pitcairn]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [PL] , Country or Area Name : [Poland]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [PT] , Country or Area Name : [Portugal]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [PR] , Country or Area Name : [Puerto Rico]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [QA] , Country or Area Name : [Qatar]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [BL] , Country or Area Name : [Saint-Barthélemy]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [WS] , Country or Area Name : [Samoa]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [SM] , Country or Area Name : [San Marino]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [SA] , Country or Area Name : [Saudi Arabia]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [SC] , Country or Area Name : [Seychelles]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [SG] , Country or Area Name : [Singapore]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [SK] , Country or Area Name : [Slovakia]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [SI] , Country or Area Name : [Slovenia]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [SO] , Country or Area Name : [Somalia]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [GS] , Country or Area Name : [South Georgia and the South Sandwich Islands]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [SS] , Country or Area Name : [South Sudan]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [ES] , Country or Area Name : [Spain]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [SE] , Country or Area Name : [Sweden]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [SY] , Country or Area Name : [Syrian Arab Republic (Syria)]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [TJ] , Country or Area Name : [Tajikistan]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [TZ] , Country or Area Name : [Tanzania, United Republic of]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [TH] , Country or Area Name : [Thailand]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [TL] , Country or Area Name : [Timor-Leste]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [TT] , Country or Area Name : [Trinidad and Tobago]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [TN] , Country or Area Name : [Tunisia]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [TR] , Country or Area Name : [Turkey]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [TV] , Country or Area Name : [Tuvalu]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [UG] , Country or Area Name : [Uganda]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [UA] , Country or Area Name : [Ukraine]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [AE] , Country or Area Name : [United Arab Emirates]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [US] , Country or Area Name : [United States of America]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [UM] , Country or Area Name : [US Minor Outlying Islands]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [UY] , Country or Area Name : [Uruguay]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [VN] , Country or Area Name : [Viet Nam]
2017/01/01 00:00:00 ISO ALPHA-2 Code : [YE] , Country or Area Name : [Yemen]
2017/01/01 00:00:00 the count of number of rows in the mdbm(=/tmp/example1.mdbm) is `246` rows
2017/01/01 00:00:00 complete
```

## Gets the count of number of records

### Code
```go
package main

import (
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	mdbm "github.com/torden/go-mdbm"
)

const (
	mdbmPath1   = "/tmp/example1.mdbm"
	mdbmPath2   = "/tmp/example2.mdbm"
	sample1Path = "./sample1.tsv" //ISO Alpha-2,3 and Numeric Country Codes
)

func main() {

	var err error

	log.Printf("Gets the count of number of records of the mdbm file(=%s)", mdbmPath1)

	//init. the go-mdbm
	dbm := mdbm.NewMDBM()

	//create & open(RDRW) an mdbm file
	err = dbm.EasyOpen(mdbmPath1, 0644)

	//the mdbm object close at close func
	defer dbm.EasyClose()

	//check the open error
	if err != nil {
		log.Fatalf("failed, can't open mdbm file\npath=%s, err=%v", mdbmPath1, err)
	}

	//obain the count of number of rows in the MDBM file
	cnt, err := dbm.EasyGetNumOfRows()
	if err != nil {
		log.Fatalf("failed, can't obtain num of rows in the mdbm file\npath=%s, err=%v", mdbmPath1, err)
	}

	log.Printf("the count of number of rows in the mdbm(=%s) is `%d` rows", mdbmPath1, cnt)
	log.Println("complete")

}
```


### Output
```
Press ENTER or type command to continue
2017/12/19 16:03:54 Gets the count of number of records of the mdbm file(=/tmp/example1.mdbm)
2017/12/19 16:03:54 the count of number of rows in the mdbm(=/tmp/example1.mdbm) is `246` rows
2017/12/19 16:03:54 complete
```


## Gets the list of keys in-place

### Code
```go
package main

import (
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	mdbm "github.com/torden/go-mdbm"
)

const (
	mdbmPath1   = "/tmp/example1.mdbm"
	mdbmPath2   = "/tmp/example2.mdbm"
	sample1Path = "./sample1.tsv" //ISO Alpha-2,3 and Numeric Country Codes
)

func main() {

	var err error

	log.Printf("Gets the list of of key of the mdbm file(=%s)", mdbmPath1)

	//init. the go-mdbm
	dbm := mdbm.NewMDBM()

	//create & open(RDRW) an mdbm file
	err = dbm.EasyOpen(mdbmPath1, 0644)

	//the mdbm object close at close func
	defer dbm.EasyClose()

	//check the open error
	if err != nil {
		log.Fatalf("failed, can't open mdbm file\npath=%s, err=%v", mdbmPath1, err)
	}

	//obtain the list of key in the MDBM file
	keyList, err := dbm.EasyGetKeyList()
	if err != nil {
		log.Fatalf("failed, can't obtain the list of key in the mdbm file\npath=%s, err=%v", mdbmPath1, err)
	}

	for k, v := range keyList {
		log.Printf("[%d] %s\n", k, v)
	}

	log.Println("complete")
}
```


### Runs the Random Deleting records 

```go
package main

import (
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	mdbm "github.com/torden/go-mdbm"
)

const (
	mdbmPath1   = "/tmp/example1.mdbm"
	mdbmPath2   = "/tmp/example2.mdbm"
	sample1Path = "./sample1.tsv" //ISO Alpha-2,3 and Numeric Country Codes
)

func main() {

	var rv int
	var err error
	err = os.Remove(mdbmPath1)
	if err == nil {
		log.Println("not exists the mdbm file(=%s)", mdbmPath1)
	} else {
		log.Println("remove the mdbm file(=%s)", mdbmPath1)
	}
	err = os.Remove(mdbmPath2)
		if err == nil {
		log.Println("not exists the mdbm file(=%s)", mdbmPath2)
	} else {
		log.Println("remove the mdbm file(=%s)", mdbmPath2)
	}

	log.Printf("Runs the random delete records of the mdbm file(=%s)", mdbmPath1)

	//init. the go-mdbm
	dbm := mdbm.NewMDBM()

	//create & open(RDRW) an mdbm file
	err = dbm.EasyOpen(mdbmPath1, 0644)

	//the mdbm object close at close func
	defer dbm.EasyClose()

	//check the open error
	if err != nil {
		log.Fatalf("failed, can't open mdbm file\npath=%s, err=%v", mdbmPath1, err)
	}

	//obtain the list of key in the MDBM file
	keyList, err := dbm.EasyGetKeyList()
	if err != nil {
		log.Fatalf("failed, can't obtain the list of key in the mdbm file\npath=%s, err=%v", mdbmPath1, err)
	}

	keySize := len(keyList)

	//obtain a pseudo-random 32-bit value as a uint32
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	var rkey int

	for i := 0; i <= cnt; i++ {

		//obtain a pseudo-random 32-bit value between 0 and keySize
		for {
			rkey = random.Intn(keySize)
			if rkey >= 0 {
				break
			}
		}

		//delete a row
		_, err := dbm.DeleteWithLock(keyList[rkey])
		if err != nil { // || !rv {
			log.Printf("not exists key(=%s)", keyList[rkey])
		} else {
			log.Printf("deleted key(=%s)", keyList[rkey])
		}
	}
}
```

### Output
```
2017/12/19 16:04:32 Gets the list of of key of the mdbm file(=/tmp/example1.mdbm)
2017/12/19 16:04:32 [0] AF
2017/12/19 16:04:32 [1] AX
2017/12/19 16:04:32 [2] AL
2017/12/19 16:04:32 [3] AD
2017/12/19 16:04:32 [4] AR
2017/12/19 16:04:32 [5] AT
2017/12/19 16:04:32 [6] AZ
2017/12/19 16:04:32 [7] BS
2017/12/19 16:04:32 [8] BY
2017/12/19 16:04:32 [9] BE
2017/12/19 16:04:32 [10] BM
2017/12/19 16:04:32 [11] BO
2017/12/19 16:04:32 [12] BA
2017/12/19 16:04:32 [13] BW
2017/12/19 16:04:32 [14] VG
2017/12/19 16:04:32 [15] BG
2017/12/19 16:04:32 [16] ZA
2017/12/19 16:04:32 [17] BI
2017/12/19 16:04:32 [18] KH
2017/12/19 16:04:32 [19] CV
2017/12/19 16:04:32 [20] CF
2017/12/19 16:04:32 [21] CL
2017/12/19 16:04:32 [22] CN
2017/12/19 16:04:32 [23] HK
2017/12/19 16:04:32 [24] CX
2017/12/19 16:04:32 [25] CD
2017/12/19 16:04:32 [26] CR
2017/12/19 16:04:32 [27] CZ
2017/12/19 16:04:32 [28] DK
2017/12/19 16:04:32 [29] SR
2017/12/19 16:04:32 [30] DM
2017/12/19 16:04:32 [31] DO
2017/12/19 16:04:32 [32] SV
2017/12/19 16:04:32 [33] ER
2017/12/19 16:04:32 [34] ET
2017/12/19 16:04:32 [35] FK
2017/12/19 16:04:32 [36] FO
2017/12/19 16:04:32 [37] FI
2017/12/19 16:04:32 [38] GF
2017/12/19 16:04:32 [39] DE
2017/12/19 16:04:32 [40] GH
2017/12/19 16:04:32 [41] GR
2017/12/19 16:04:32 [42] GL
2017/12/19 16:04:32 [43] GD
2017/12/19 16:04:32 [44] GP
2017/12/19 16:04:32 [45] GT
2017/12/19 16:04:32 [46] GN
2017/12/19 16:04:32 [47] HM
2017/12/19 16:04:32 [48] VA
2017/12/19 16:04:32 [49] HU
2017/12/19 16:04:32 [50] IN
2017/12/19 16:04:32 [51] ID
2017/12/19 16:04:32 [52] IR
2017/12/19 16:04:32 [53] IL
2017/12/19 16:04:32 [54] IT
2017/12/19 16:04:32 [55] JM
2017/12/19 16:04:32 [56] SD
2017/12/19 16:04:32 [57] JE
2017/12/19 16:04:32 [58] JO
2017/12/19 16:04:32 [59] KZ
2017/12/19 16:04:32 [60] KP
2017/12/19 16:04:32 [61] KR
2017/12/19 16:04:32 [62] LA
2017/12/19 16:04:32 [63] LS
2017/12/19 16:04:32 [64] LY
2017/12/19 16:04:32 [65] LI
2017/12/19 16:04:32 [66] LU
2017/12/19 16:04:32 [67] MV
2017/12/19 16:04:32 [68] ML
2017/12/19 16:04:32 [69] MT
2017/12/19 16:04:32 [70] MH
2017/12/19 16:04:32 [71] MR
2017/12/19 16:04:32 [72] YT
2017/12/19 16:04:32 [73] MX
2017/12/19 16:04:32 [74] FM
2017/12/19 16:04:32 [75] MD
2017/12/19 16:04:32 [76] MN
2017/12/19 16:04:32 [77] MZ
2017/12/19 16:04:32 [78] NA
2017/12/19 16:04:32 [79] AN
2017/12/19 16:04:32 [80] NC
2017/12/19 16:04:32 [81] NI
2017/12/19 16:04:32 [82] NE
2017/12/19 16:04:32 [83] NG
2017/12/19 16:04:32 [84] NU
2017/12/19 16:04:32 [85] MP
2017/12/19 16:04:32 [86] NO
2017/12/19 16:04:32 [87] PK
2017/12/19 16:04:32 [88] PW
2017/12/19 16:04:32 [89] PS
2017/12/19 16:04:32 [90] PA
2017/12/19 16:04:32 [91] PG
2017/12/19 16:04:32 [92] PY
2017/12/19 16:04:32 [93] PE
2017/12/19 16:04:32 [94] RE
2017/12/19 16:04:32 [95] RO
2017/12/19 16:04:32 [96] RU
2017/12/19 16:04:32 [97] RW
2017/12/19 16:04:32 [98] SH
2017/12/19 16:04:32 [99] KN
2017/12/19 16:04:32 [100] LC
2017/12/19 16:04:32 [101] MF
2017/12/19 16:04:32 [102] PM
2017/12/19 16:04:32 [103] VC
2017/12/19 16:04:32 [104] ST
2017/12/19 16:04:32 [105] SN
2017/12/19 16:04:32 [106] RS
2017/12/19 16:04:32 [107] SL
2017/12/19 16:04:32 [108] SB
2017/12/19 16:04:32 [109] LK
2017/12/19 16:04:32 [110] SJ
2017/12/19 16:04:32 [111] SZ
2017/12/19 16:04:32 [112] CH
2017/12/19 16:04:32 [113] TW
2017/12/19 16:04:32 [114] TG
2017/12/19 16:04:32 [115] TK
2017/12/19 16:04:32 [116] TO
2017/12/19 16:04:32 [117] TM
2017/12/19 16:04:32 [118] TC
2017/12/19 16:04:32 [119] GB
2017/12/19 16:04:32 [120] UZ
2017/12/19 16:04:32 [121] VU
2017/12/19 16:04:32 [122] VE
2017/12/19 16:04:32 [123] VI
2017/12/19 16:04:32 [124] WF
2017/12/19 16:04:32 [125] EH
2017/12/19 16:04:32 [126] ZM
2017/12/19 16:04:32 [127] ZW
2017/12/19 16:04:32 [128] DZ
2017/12/19 16:04:32 [129] AS
2017/12/19 16:04:32 [130] AO
2017/12/19 16:04:32 [131] AI
2017/12/19 16:04:32 [132] AQ
2017/12/19 16:04:32 [133] AG
2017/12/19 16:04:32 [134] AM
2017/12/19 16:04:32 [135] AW
2017/12/19 16:04:32 [136] AU
2017/12/19 16:04:32 [137] BH
2017/12/19 16:04:32 [138] BD
2017/12/19 16:04:32 [139] BB
2017/12/19 16:04:32 [140] BZ
2017/12/19 16:04:32 [141] BJ
2017/12/19 16:04:32 [142] BT
2017/12/19 16:04:32 [143] BV
2017/12/19 16:04:32 [144] BR
2017/12/19 16:04:32 [145] IO
2017/12/19 16:04:32 [146] BN
2017/12/19 16:04:32 [147] BF
2017/12/19 16:04:32 [148] CM
2017/12/19 16:04:32 [149] CA
2017/12/19 16:04:32 [150] KY
2017/12/19 16:04:32 [151] TD
2017/12/19 16:04:32 [152] MO
2017/12/19 16:04:32 [153] CC
2017/12/19 16:04:32 [154] CO
2017/12/19 16:04:32 [155] KM
2017/12/19 16:04:32 [156] CG
2017/12/19 16:04:32 [157] CK
2017/12/19 16:04:32 [158] CI
2017/12/19 16:04:32 [159] HR
2017/12/19 16:04:32 [160] CU
2017/12/19 16:04:32 [161] CY
2017/12/19 16:04:32 [162] DJ
2017/12/19 16:04:32 [163] EC
2017/12/19 16:04:32 [164] EG
2017/12/19 16:04:32 [165] GQ
2017/12/19 16:04:32 [166] EE
2017/12/19 16:04:32 [167] FJ
2017/12/19 16:04:32 [168] FR
2017/12/19 16:04:32 [169] PF
2017/12/19 16:04:32 [170] TF
2017/12/19 16:04:32 [171] GA
2017/12/19 16:04:32 [172] GM
2017/12/19 16:04:32 [173] GE
2017/12/19 16:04:32 [174] GI
2017/12/19 16:04:32 [175] GU
2017/12/19 16:04:32 [176] GG
2017/12/19 16:04:32 [177] GW
2017/12/19 16:04:32 [178] GY
2017/12/19 16:04:32 [179] HT
2017/12/19 16:04:32 [180] HN
2017/12/19 16:04:32 [181] IS
2017/12/19 16:04:32 [182] IQ
2017/12/19 16:04:32 [183] IE
2017/12/19 16:04:32 [184] IM
2017/12/19 16:04:32 [185] JP
2017/12/19 16:04:32 [186] KE
2017/12/19 16:04:32 [187] KI
2017/12/19 16:04:32 [188] KW
2017/12/19 16:04:32 [189] KG
2017/12/19 16:04:32 [190] LV
2017/12/19 16:04:32 [191] LB
2017/12/19 16:04:32 [192] LR
2017/12/19 16:04:32 [193] LT
2017/12/19 16:04:32 [194] MK
2017/12/19 16:04:32 [195] MG
2017/12/19 16:04:32 [196] MW
2017/12/19 16:04:32 [197] MY
2017/12/19 16:04:32 [198] MQ
2017/12/19 16:04:32 [199] MU
2017/12/19 16:04:32 [200] MC
2017/12/19 16:04:32 [201] ME
2017/12/19 16:04:32 [202] MS
2017/12/19 16:04:32 [203] MA
2017/12/19 16:04:32 [204] MM
2017/12/19 16:04:32 [205] NR
2017/12/19 16:04:32 [206] NP
2017/12/19 16:04:32 [207] NL
2017/12/19 16:04:32 [208] NZ
2017/12/19 16:04:32 [209] NF
2017/12/19 16:04:32 [210] OM
2017/12/19 16:04:32 [211] PH
2017/12/19 16:04:32 [212] PN
2017/12/19 16:04:32 [213] PL
2017/12/19 16:04:32 [214] PT
2017/12/19 16:04:32 [215] PR
2017/12/19 16:04:32 [216] QA
2017/12/19 16:04:32 [217] BL
2017/12/19 16:04:32 [218] WS
2017/12/19 16:04:32 [219] SM
2017/12/19 16:04:32 [220] SA
2017/12/19 16:04:32 [221] SC
2017/12/19 16:04:32 [222] SG
2017/12/19 16:04:32 [223] SK
2017/12/19 16:04:32 [224] SI
2017/12/19 16:04:32 [225] SO
2017/12/19 16:04:32 [226] GS
2017/12/19 16:04:32 [227] SS
2017/12/19 16:04:32 [228] ES
2017/12/19 16:04:32 [229] SE
2017/12/19 16:04:32 [230] SY
2017/12/19 16:04:32 [231] TJ
2017/12/19 16:04:32 [232] TZ
2017/12/19 16:04:32 [233] TH
2017/12/19 16:04:32 [234] TL
2017/12/19 16:04:32 [235] TT
2017/12/19 16:04:32 [236] TN
2017/12/19 16:04:32 [237] TR
2017/12/19 16:04:32 [238] TV
2017/12/19 16:04:32 [239] UG
2017/12/19 16:04:32 [240] UA
2017/12/19 16:04:32 [241] AE
2017/12/19 16:04:32 [242] US
2017/12/19 16:04:32 [243] UM
2017/12/19 16:04:32 [244] UY
2017/12/19 16:04:32 [245] VN
2017/12/19 16:04:32 [246] YE
2017/12/19 16:04:32 complete
```

## Updating records in-place

```go
package main

import (
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	mdbm "github.com/torden/go-mdbm"
)

const (
	mdbmPath1   = "/tmp/example1.mdbm"
	mdbmPath2   = "/tmp/example2.mdbm"
	sample1Path = "./sample1.tsv" //ISO Alpha-2,3 and Numeric Country Codes
)

func main() {

	var err error
	
	log.Printf("Runs the updating records of the mdbm file(=%s)", mdbmPath1)

	//init. the go-mdbm
	dbm := mdbm.NewMDBM()

	//create & open(RDRW) an mdbm file
	err = dbm.EasyOpen(mdbmPath1, 0644)

	//the mdbm object close at close func
	defer dbm.EasyClose()

	//check the open error
	if err != nil {
		log.Fatalf("failed, can't open mdbm file\npath=%s, err=%v", mdbmPath1, err)
	}

	//read a content
	data, err := ioutil.ReadFile(sample1Path)
	if err != nil {
		log.Fatalf("failed, read to the tsv file\npath=%s, err=%v", sample1Path, err)
	}

	//convert []byte to string and split by newline
	dataStrAr := strings.Split(string(data), "\n")
	for k, v := range dataStrAr {

		if k == 0 { //header
			continue
		}

		//obtain data by field
		row := strings.Split(v, "\t")
		if len(row) < 4 {
			continue
		}

		//0 : Country or Area Name
		//1 : ISO ALPHA-2 Code
		//2 : ISO ALPHA-3 Code
		//3 : USO Numeric Code , UN M49 Numeric Code
		rv, err := dbm.StoreWithLock(row[0], row[3], mdbm.Replace) // if key does exist; update to value
		if err != nil {
			log.Fatalf("failed, can't data(key=%+v, value=%+v) add to the mdbm file(=%s)\nrv=%d, err=%v", row[0], row[1], mdbmPath1, rv, err)
		}
	}

	log.Println("complete")
}
```

### Output
```
2017/01/01 00:00:00 Runs the updating records of the mdbm file(=/tmp/example1.mdbm)
2017/01/01 00:00:00 complete
```


## Todo

* coverage up to 95% (min.)
* parallel testing..
* data race detecting testing..
* Binding All APIs without deprecated apis
* Testing on another platform (osx, bsd...)
* Pre-compiled mdbm library by OS
* Stabilization

---
Please feel free.
