package main

import (
	"log"

	mdbm "github.com/torden/go-mdbm"
)

const dbPath = "/tmp/test_benchmark_large.mdbm"

func main() {
	dbm := mdbm.NewMDBM()
	err := dbm.Open("/tmp/test_benchmark_large.mdbm", mdbm.Create|mdbm.Rdrw|mdbm.LargeObjects|mdbm.Trunc, 0644, 0, 0)
	defer dbm.EasyClose()
	if err != nil {
		log.Fatalf("failed, can't open mdbm file\npath=%s, err=%v", dbm.GetDBMFile(), err)
	}

	for i := uint32(0); i <= 100000000; i++ {

		rv, err := dbm.Store(i, i, mdbm.Insert)
		if err != nil {
			log.Fatalf("failed, can't data(kv=%+v) add to the mdbm file(=%s)\nrv=%d, err=%v", i, dbm.GetDBMFile(), rv, err)
		}
	}

	log.Println("complete")
}
