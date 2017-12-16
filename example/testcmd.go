package main

import (
	"io/ioutil"
	"log"
	"math/rand"
	"strings"
	"time"

	mdbm "github.com/torden/go-mdbm"
)

const (
	mdbmPath    = "/tmp/example1.mdbm"
	sample1Path = "./sample1.tsv" //ISO Alpha-2,3 and Numeric Country Codes
)

func exampleGenerateMDBMFile() {

	log.Printf("generating the mdbm file(=%s)", mdbmPath)

	//init. the go-mdbm
	dbm := mdbm.NewMDBM()

	//create & open(RDRW) an mdbm file
	err := dbm.EasyOpen(mdbmPath, 0644)

	//the mdbm object close at close func
	defer dbm.EasyClose()

	//check the open error
	if err != nil {
		log.Fatalf("failed, can't open mdbm file\npath=%s, err=%v", mdbmPath, err)
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
			log.Fatalf("failed, can't data(key=%+v, value=%+v) add to the mdbm file(=%s)\nrv=%d, err=%v", row[0], row[1], mdbmPath, rv, err)
		}
	}

	log.Println("complete")
}

func exampleFetch(limit int) {

	//init. the go-mdbm
	dbm := mdbm.NewMDBM()

	//create & open(RDRW) an mdbm file
	err := dbm.EasyOpen(mdbmPath, 0644)

	//the mdbm object close at close func
	defer dbm.EasyClose()

	//check the open error
	if err != nil {
		log.Fatalf("failed, can't open mdbm file\npath=%s, err=%v", mdbmPath, err)
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
			log.Fatalf("failed, can't find out value in the mdbm file(=%s)\nekey=%s, err=%v", mdbmPath, keyList[rkey], err)
		}

		log.Printf("ISO ALPHA-2 Code of [%s] is [%s]", keyList[rkey], val)
	}

	log.Println("complete")
}

func exampleIterationUseFirstNext() {

	//init. the go-mdbm
	dbm := mdbm.NewMDBM()

	//create & open(RDRW) an mdbm file
	err := dbm.EasyOpen(mdbmPath, 0644)

	//the mdbm object close at close func
	defer dbm.EasyClose()

	//check the open error
	if err != nil {
		log.Fatalf("failed, can't open mdbm file\npath=%s, err=%v", mdbmPath, err)
	}

	key, val, err := dbm.First()
	if err != nil {
		log.Fatalf("key=%s, val=%s, err=%v", key, val, err)
	}

	log.Printf("ISO ALPHA-2 Code of [%s] is [%s]", key, val)

	var cnt uint64

	for {

		key, val, err := dbm.Next()
		if len(key) < 1 {
			break
		}

		if err != nil {
			log.Fatalf("key=%s, val=%s, err=%v", key, val, err)
		}

		log.Printf("ISO ALPHA-2 Code of [%s] is [%s]", key, val)
		cnt++
	}

	log.Printf("the count of number of rows in the mdbm(=%s) is `%d` rows", mdbmPath, cnt)
}

func exampleNumRows() {

	//init. the go-mdbm
	dbm := mdbm.NewMDBM()

	//create & open(RDRW) an mdbm file
	err := dbm.EasyOpen(mdbmPath, 0644)

	//the mdbm object close at close func
	defer dbm.EasyClose()

	//check the open error
	if err != nil {
		log.Fatalf("failed, can't open mdbm file\npath=%s, err=%v", mdbmPath, err)
	}

	//obain the count of number of rows in the MDBM file
	cnt, err := dbm.EasyGetNumOfRows()
	if err != nil {
		log.Fatalf("failed, can't obtain num of rows in the mdbm file\npath=%s, err=%v", mdbmPath, err)
	}

	log.Printf("the count of number of rows in the mdbm(=%s) is `%d` rows", mdbmPath, cnt)

}

func exampleKeyList() {

	//init. the go-mdbm
	dbm := mdbm.NewMDBM()

	//create & open(RDRW) an mdbm file
	err := dbm.EasyOpen(mdbmPath, 0644)

	//the mdbm object close at close func
	defer dbm.EasyClose()

	//check the open error
	if err != nil {
		log.Fatalf("failed, can't open mdbm file\npath=%s, err=%v", mdbmPath, err)
	}

	//obtain the list of key in the MDBM file
	keyList, err := dbm.EasyGetKeyList()
	if err != nil {
		log.Fatalf("failed, can't obtain the list of key in the mdbm file\npath=%s, err=%v", mdbmPath, err)
	}

	for k, v := range keyList {
		log.Printf("[%d] %s\n", k, v)
	}
}

func main() {

	//os.Remove(mdbmPath)
	exampleGenerateMDBMFile()
	exampleFetch(10)
	exampleIterationUseFirstNext()
	exampleNumRows()
	exampleKeyList()
}
