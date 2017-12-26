package mdbm_test

import (
	"log"
	"testing"
	"time"

	mdbm "github.com/torden/go-mdbm"
)

func Benchmark_mdbm_Store(b *testing.B) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBMBenchmark1, 0644)
	defer dbm.EasyClose()
	if err != nil {
		log.Fatalf("failured, can't open the mdbm, path=%s, err=%v", dbm.GetDBMFile(), err)
	}

	for i := 0; i < b.N; i++ {
		rv, err := dbm.Store(i, time.Now().UnixNano(), mdbm.Replace)
		if err != nil {
			if err != nil {
				log.Fatalf("failed, can't data(=%d) add to the mdbm file(=%s), rv=%d, err=%v", i, dbm.GetDBMFile(), rv, err)
			}
		}
	}
}

func Benchmark_mdbm_StoreWithLock(b *testing.B) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBMBenchmark2, 0644)
	defer dbm.EasyClose()
	if err != nil {
		log.Fatalf("failured, can't open the mdbm, path=%s, err=%v", dbm.GetDBMFile(), err)
	}

	for i := 0; i < b.N; i++ {
		rv, err := dbm.StoreWithLock(i, time.Now().UnixNano(), mdbm.Replace)
		if err != nil {
			if err != nil {
				log.Fatalf("failed, can't data(=%d) add to the mdbm file(=%s), rv=%d, err=%v", i, dbm.GetDBMFile(), rv, err)
			}
		}
	}
}

func Benchmark_mdbm_StoreOnLock(b *testing.B) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBMBenchmark2, 0644)
	defer dbm.EasyClose()
	if err != nil {
		log.Fatalf("failured, can't open the mdbm, path=%s, err=%v", dbm.GetDBMFile(), err)
	}

	rv, err := dbm.Lock()
	if err != nil {
		log.Fatalf("failed, can't obtain lock, path=%s, rv=%d, err=%v", dbm.GetDBMFile(), rv, err)
	}

	defer dbm.Unlock()

	for i := 0; i < b.N; i++ {
		rv, err := dbm.Store(i, time.Now().UnixNano(), mdbm.Replace)
		if err != nil {
			if err != nil {
				log.Fatalf("failed, can't data(=%d) add to the mdbm file(=%s), rv=%d, err=%v", i, dbm.GetDBMFile(), rv, err)
			}
		}
	}
}

func Benchmark_mdbm_Fetch(b *testing.B) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBMBenchmark1, 0644)
	defer dbm.EasyClose()
	if err != nil {
		log.Fatalf("failured, can't open the mdbm, path=%s, err=%v", dbm.GetDBMFile(), err)
	}

	for i := 0; i < b.N; i++ {

		_, err = dbm.Fetch(i)
		if err != nil {
			continue
		}
	}
}

func Benchmark_mdbm_FetchWithLock(b *testing.B) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBMBenchmark1, 0644)
	defer dbm.EasyClose()
	if err != nil {
		log.Fatalf("failured, can't open the mdbm, path=%s, err=%v", dbm.GetDBMFile(), err)
	}

	for i := 0; i < b.N; i++ {

		_, err = dbm.FetchWithLock(i)
		if err != nil {
			continue
		}
	}
}

func Benchmark_mdbm_FetchOnLock(b *testing.B) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBMBenchmark1, 0644)
	defer dbm.EasyClose()
	if err != nil {
		log.Fatalf("failured, can't open the mdbm, path=%s, err=%v", dbm.GetDBMFile(), err)
	}

	rv, err := dbm.Lock()
	if err != nil {
		log.Fatalf("failed, can't obtain lock, path=%s, rv=%d, err=%v", dbm.GetDBMFile(), rv, err)
	}

	defer dbm.Unlock()

	for i := 0; i < b.N; i++ {

		_, err = dbm.FetchWithLock(i)
		if err != nil {
			continue
		}
	}
}

func Benchmark_mdbm_PreLoad_Fetch(b *testing.B) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBMBenchmark1, 0644)
	defer dbm.EasyClose()
	if err != nil {
		log.Fatalf("failured, can't open the mdbm, path=%s, err=%v", dbm.GetDBMFile(), err)
	}

	rv, err := dbm.PreLoad()
	if err != nil {
		log.Fatalf("failured, can't pre-load the mdbm, path=%s, rv=%d, err=%v", dbm.GetDBMFile(), rv, err)
	}

	for i := 0; i < b.N; i++ {

		_, err = dbm.Fetch(i)
		if err != nil {
			continue
		}
	}
}

func Benchmark_mdbm_PreLoad_FetchWithLock(b *testing.B) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBMBenchmark1, 0644)
	defer dbm.EasyClose()
	if err != nil {
		log.Fatalf("failured, can't open the mdbm, path=%s, err=%v", dbm.GetDBMFile(), err)
	}

	rv, err := dbm.PreLoad()
	if err != nil {
		log.Fatalf("failured, can't pre-load the mdbm, path=%s, rv=%d, err=%v", dbm.GetDBMFile(), rv, err)
	}

	for i := 0; i < b.N; i++ {

		_, err = dbm.FetchWithLock(i)
		if err != nil {
			continue
		}
	}
}

func Benchmark_mdbm_PreLoad_FetchOnLock(b *testing.B) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBMBenchmark1, 0644)
	defer dbm.EasyClose()
	if err != nil {
		log.Fatalf("failured, can't open the mdbm, path=%s, err=%v", dbm.GetDBMFile(), err)
	}

	rv, err := dbm.PreLoad()
	if err != nil {
		log.Fatalf("failured, can't pre-load the mdbm, path=%s, rv=%d, err=%v", dbm.GetDBMFile(), rv, err)
	}

	rv, err = dbm.Lock()
	if err != nil {
		log.Fatalf("failed, can't obtain lock, path=%s, rv=%d, err=%v", dbm.GetDBMFile(), rv, err)
	}

	defer dbm.Unlock()

	for i := 0; i < b.N; i++ {

		_, err = dbm.FetchWithLock(i)
		if err != nil {
			continue
		}
	}
}
