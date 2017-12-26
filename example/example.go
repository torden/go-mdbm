package main

import (
	"fmt"
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
	mdbmPathDup = "/tmp/exampleDup.mdbm"
	sample1Path = "./sample1.tsv" //ISO Alpha-2,3 and Numeric Country Codes
)

func exampleGenerateMDBMFile() {

	log.Printf("Creating and Populating the mdbm file(=%s)", mdbmPath1)

	//init. the go-mdbm
	dbm := mdbm.NewMDBM()

	//create & open(RDRW) an mdbm file
	err := dbm.EasyOpen(mdbmPath1, 0644)

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

		//0 : ISO ALPHA-2 Code
		//1 : Country or Area Name
		//2 : ISO ALPHA-3 Code
		//3 : USO Numeric Code , UN M49 Numeric Code

		rv, err := dbm.Store(row[0], row[1], mdbm.Insert) // if key does not exist; fail if exists
		if err != nil {
			log.Fatalf("failed, can't data(key=%+v, value=%+v) add to the mdbm file(=%s)\nrv=%d, err=%v", row[0], row[1], mdbmPath1, rv, err)
		}
	}

	log.Println("complete")
}

func exampleDupDataMDBMFile() {

	log.Printf("Creating and Populating the mdbm file(=%s)", mdbmPath1)

	var rv int
	var err error

	//init. the go-mdbm
	dbm := mdbm.NewMDBM()

	//create & open(RDRW|DUP) an mdbm file
	err = dbm.Open(mdbmPathDup, mdbm.Create|mdbm.Rdrw|mdbm.LargeObjects|mdbm.InsertDup, 0644, 0, 0)

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

		//0 : ISO ALPHA-2 Code
		//1 : Country or Area Name
		//2 : ISO ALPHA-3 Code
		//3 : USO Numeric Code , UN M49 Numeric Code

		rv, err = dbm.Store(row[0], row[1], mdbm.InsertDup)
		if err != nil {
			log.Fatalf("failed, can't data(key=%+v, value=%+v) add to the mdbm file(=%s)\nrv=%d, err=%v", row[0], row[1], mdbmPath1, rv, err)
		}

		rv, err = dbm.Store(row[0], row[2], mdbm.InsertDup)
		if err != nil {
			log.Fatalf("failed, can't data(key=%+v, value=%+v) add to the mdbm file(=%s)\nrv=%d, err=%v", row[0], row[1], mdbmPath1, rv, err)
		}

		rv, err = dbm.Store(row[0], row[3], mdbm.InsertDup)
		if err != nil {
			log.Fatalf("failed, can't data(key=%+v, value=%+v) add to the mdbm file(=%s)\nrv=%d, err=%v", row[0], row[1], mdbmPath1, rv, err)
		}

	}

	log.Println("complete")
}

func exampleGenerateUseAnyDataTypeMDBMFile() {

	var rv int
	var err error

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
		if len(key) < 1 {
			break
		}

		if err != nil {
			log.Fatalf("key=%s, val=%s, err=%v", key, val, err)
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

func exampleFetch(limit int) {

	log.Printf("Fetching records in the mdbm file(=%s)", mdbmPath1)

	//init. the go-mdbm
	dbm := mdbm.NewMDBM()

	//create & open(RDRW) an mdbm file
	err := dbm.EasyOpen(mdbmPath1, 0644)

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
	for i := 0; i <= limit; i++ {

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

func exampleFetchDup(limit int) {

	log.Printf("Fetching records in the mdbm file(=%s)", mdbmPath1)

	//init. the go-mdbm
	dbm := mdbm.NewMDBM()

	//create & open(ReadOnly) an mdbm file
	err := dbm.Open(mdbmPathDup, mdbm.Rdonly, 0644, 0, 0)

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

	//header
	fmt.Println(strings.Repeat("-", 129))
	fmt.Printf("| %-30s| %-30s| %-30s| %-30s|\n", "ISO ALPHA-2 Code", "Country or Area Name", "ISO ALPHA-3 Code", "USO Numeric Code")
	fmt.Println(strings.Repeat("-", 129))

	//random fetching..
	for i := 0; i <= limit; i++ {

		//obtain a pseudo-random 32-bit value between 0 and keySize
		for {
			rkey = random.Intn(keySize)
			if rkey >= 0 {
				break
			}
		}

		fmt.Printf("| %-30s", keyList[rkey])

		rv := 0
		var val string
		iter := dbm.GetNewIter()
		for rv != -1 {

			rv, val, _, err = dbm.FetchDupRWithLock(keyList[rkey], &iter)
			if rv == -1 {
				break
			}
			fmt.Printf("| %-30s", val)
		}

		fmt.Println("|")
	}

	//footer
	fmt.Println(strings.Repeat("-", 129))

	log.Println("complete")
}

func exampleIterationUseFirstNext() {

	log.Printf("Iterating over all records in the mdbm file(=%s)", mdbmPath1)

	//init. the go-mdbm
	dbm := mdbm.NewMDBM()

	//create & open(RDRW) an mdbm file
	err := dbm.EasyOpen(mdbmPath1, 0644)

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

	log.Printf("ISO ALPHA-2 Code : [%s] , Country or Area Name : [%s]", key, val)

	var cnt uint64

	for {

		key, val, err := dbm.Next()
		if len(key) < 1 {
			break
		}

		if err != nil {
			log.Fatalf("key=%s, val=%s, err=%v", key, val, err)
		}

		log.Printf("ISO ALPHA-2 Code : [%s] , Country or Area Name : [%s]", key, val)
		cnt++
	}

	log.Printf("the count of number of rows in the mdbm(=%s) is `%d` rows", mdbmPath1, cnt)
	log.Println("complete")
}

func exampleNumRows() {

	log.Printf("Gets the count of number of records of the mdbm file(=%s)", mdbmPath1)

	//init. the go-mdbm
	dbm := mdbm.NewMDBM()

	//create & open(RDRW) an mdbm file
	err := dbm.EasyOpen(mdbmPath1, 0644)

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

func exampleKeyList() {

	log.Printf("Gets the list of of key of the mdbm file(=%s)", mdbmPath1)

	//init. the go-mdbm
	dbm := mdbm.NewMDBM()

	//create & open(RDRW) an mdbm file
	err := dbm.EasyOpen(mdbmPath1, 0644)

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

func exampleDelete(cnt int) {

	log.Printf("Runs the random delete records of the mdbm file(=%s)", mdbmPath1)

	//init. the go-mdbm
	dbm := mdbm.NewMDBM()

	//create & open(RDRW) an mdbm file
	err := dbm.EasyOpen(mdbmPath1, 0644)

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

func exampleUpdateValue() {

	log.Printf("Runs the updating records of the mdbm file(=%s)", mdbmPath1)

	//init. the go-mdbm
	dbm := mdbm.NewMDBM()

	//create & open(RDRW) an mdbm file
	err := dbm.EasyOpen(mdbmPath1, 0644)

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

func main() {

	os.Remove(mdbmPath1)
	os.Remove(mdbmPath2)
	os.Remove(mdbmPathDup)

	exampleGenerateMDBMFile()
	exampleDupDataMDBMFile()
	exampleGenerateUseAnyDataTypeMDBMFile()
	exampleFetch(10)
	exampleFetchDup(20)
	exampleIterationUseFirstNext()
	exampleNumRows()
	exampleKeyList()
	exampleDelete(5)

	exampleUpdateValue()
	exampleIterationUseFirstNext()
}
