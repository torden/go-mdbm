package mdbm_test

import (
	"log"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/torden/go-mdbm"
	"github.com/torden/go-strutil"
)

var assert = strutils.NewAssert()

func TestMain(t *testing.T) {

	pathList := [...]string{
		pathTestDBM1, pathTestDBM2, pathTestDBM3,
		pathTestDBMLarge, pathTestDBMHash, pathTestDBMDup,
		pathTestDBMCache, pathTestDBMV2, pathTestDBMLock1,
		pathTestDBMDelete, pathTestDBMLock2, pathTestDBMAnyDataType,
		pathTestDBMStr}

	dbm := mdbm.NewMDBM()

	for _, path := range pathList {

		_, err := dbm.DeleteLockFiles(path)
		if err == nil {
			log.Printf("delete lock files of %s", path)
		}

		err = os.Remove(path)
		if err != nil {
			log.Printf("not exists the `%s` file", path)
		} else {
			log.Printf("remove the `%s` file", path)
		}

	}
}

func Test_mdbm_EasyOpen_EasyClose(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()

	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)
}

func Test_mdbm_SetCacheMode(t *testing.T) {

	var rv int
	var err error

	dbm := mdbm.NewMDBM()
	err = dbm.EasyOpen(pathTestDBMCache, 0644)
	defer dbm.EasyClose()
	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	rv, err = dbm.SetCacheMode(mdbm.CacheModeNone)
	assert.AssertNil(t, err, "failured, can't set cachemode(=mdbm.CacheModeNone), path=%s, rv=%d, err=%v", pathTestDBMCache, rv, err)

	rv, err = dbm.SetCacheMode(mdbm.CacheModeLFU)
	assert.AssertNil(t, err, "failured, can't set cachemode(=mdbm.CacheModeLfu), path=%s, rv=%d, err=%v", pathTestDBMCache, rv, err)

	rv, err = dbm.SetCacheMode(mdbm.CacheModeLRU)
	assert.AssertNil(t, err, "failured, can't set cachemode(=%s), path=%s, rv=%d, err=%v", pathTestDBMCache, rv, err)

	rv, err = dbm.SetCacheMode(mdbm.CacheModeGDSF)
	assert.AssertNil(t, err, "failured, can't set cachemode(=%s), path=%s, rv=%d, err=%v", pathTestDBMCache, rv, err)

	rv, err = dbm.SetCacheMode(mdbm.CacheModeMax)
	assert.AssertNil(t, err, "failured, can't set cachemode(=%s), path=%s, rv=%d, err=%v", pathTestDBMCache, rv, err)

	rv, err = dbm.SetCacheMode(mdbm.LargeObjects)
	assert.AssertNotNil(t, err, "failured, can't cehck wrong cachemode, path=%s, rv=%d, err=%v", pathTestDBMCache, rv, err)
}

func Test_mdbm_GetCacheModeName(t *testing.T) {

	var rv string
	var err error

	dbm := mdbm.NewMDBM()
	err = dbm.EasyOpen(pathTestDBMCache, 0644)
	defer dbm.EasyClose()
	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	rv, err = dbm.GetCacheModeName(mdbm.CacheModeNone)
	assert.AssertNil(t, err, "failured, can't get a cache mode name. err=%v\n", rv, err)
	assert.AssertEquals(t, rv, "none", "return Value mismatch.\nExpected: %v\nActual: %v", "none", rv)

	rv, err = dbm.GetCacheModeName(mdbm.CacheModeLFU)
	assert.AssertNil(t, err, "failured, can't get a cache mode name. err=%v\n", rv, err)
	assert.AssertEquals(t, rv, "LFU", "return Value mismatch.\nExpected: %v\nActual: %v", "LFU", rv)

	rv, err = dbm.GetCacheModeName(mdbm.CacheModeLRU)
	assert.AssertNil(t, err, "failured, can't get a cache mode name. err=%v\n", rv, err)
	assert.AssertEquals(t, rv, "LRU", "return Value mismatch.\nExpected: %v\nActual: %v", "LRU", rv)

	rv, err = dbm.GetCacheModeName(mdbm.CacheModeGDSF)
	assert.AssertNil(t, err, "failured, can't get a cache mode name. err=%v\n", rv, err)
	assert.AssertEquals(t, rv, "GDSF", "return Value mismatch.\nExpected: %v\nActual: %v", "GDSF", rv)

	rv, err = dbm.GetCacheModeName(mdbm.CacheModeMax)
	assert.AssertNil(t, err, "failured, can't get a cache mode name. err=%v\n", rv, err)
	assert.AssertEquals(t, rv, "GDSF", "return Value mismatch.\nExpected: %v\nActual: %v", "GDSF", rv)
}

func Test_mdbm_Open_Close(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.Open(pathTestDBM2, mdbm.Create|mdbm.Rdrw, 0644, 0, 0)
	defer dbm.Close()

	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)
}

func Test_mdbm_OrdinaryInsertData_Store1(t *testing.T) {

	var rv int
	var err error

	dbm := mdbm.NewMDBM()
	err = dbm.Open(pathTestDBM1, mdbm.Create|mdbm.Rdrw, 0644, 0, 0)
	defer dbm.Close()

	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	for i := 0; i <= loopLimit; i++ {
		rv, err = dbm.Store(i, time.Now().UnixNano(), mdbm.Replace)
		assert.AssertNil(t, err, "failured, return Value mismatch. value=%d, err=%v\n", rv, err)
	}

	rv, err = dbm.Sync()
	assert.AssertNil(t, err, "failured, execute to mdbm.Sync()\nrv=%d, err=%v", rv, err)
}

func Test_mdbm_OrdinaryInsertData_Store2(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM2, 0644)
	defer dbm.EasyClose()

	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM2, err)

	for i := 0; i <= loopLimit; i++ {
		rv, err := dbm.Store(i, time.Now().UnixNano(), mdbm.Replace)
		assert.AssertNil(t, err, "failured, return Value mismatch. value=%d, err=%v\n", rv, err)
	}
}

func Test_mdbm_OrdinaryInsertData_StoreWithLock(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM3, 0644)
	defer dbm.EasyClose()

	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM3, err)

	for i := 0; i <= loopLimit; i++ {
		rv, err := dbm.StoreWithLock(i, time.Now().UnixNano(), mdbm.Replace)
		assert.AssertNil(t, err, "failured, return Value mismatch. value=%d, err=%v\n", rv, err)
	}
}

func Test_mdbm_OrdinaryReaplceData_StoreWithLockSmart(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()

	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	for i := 0; i <= loopLimit; i++ {
		rv, err := dbm.StoreWithLockSmart(i, i, mdbm.Replace, mdbm.Rdrw)
		assert.AssertNil(t, err, "failured, return Value mismatch. value=%d, err=%v\n", rv, err)
	}
}

func Test_mdbm_OrdinaryReaplceData_StoreWithLockShared(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()

	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	for i := 0; i <= loopLimit; i++ {
		rv, err := dbm.StoreWithLockShared(i, i, mdbm.Replace)
		assert.AssertNil(t, err, "failured, return Value mismatch. value=%d, err=%v\n", rv, err)
	}
}

func Test_mdbm_OrdinaryReaplceData_StoreWithPlock(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()

	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	for i := 0; i <= loopLimit; i++ {
		rv, err := dbm.StoreWithPlock(i, i, mdbm.Replace, mdbm.Rdrw)
		assert.AssertNil(t, err, "failured, return Value mismatch. value=%d, err=%v\n", rv, err)
	}
}

func Test_mdbm_OrdinaryReaplceData_StoreWithTryLock(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()

	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	for i := 0; i <= loopLimit; i++ {
		rv, err := dbm.StoreWithTryLock(i, i, mdbm.Replace)
		assert.AssertNil(t, err, "failured, return Value mismatch. value=%d, err=%v\n", rv, err)
	}
}

func Test_mdbm_OrdinaryReaplceData_StoreWithTryLockSmart(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()

	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	for i := 0; i <= loopLimit; i++ {
		rv, err := dbm.StoreWithTryLockSmart(i, i, mdbm.Replace, mdbm.Rdrw)
		assert.AssertNil(t, err, "failured, return Value mismatch. value=%d, err=%v\n", rv, err)
	}
}

func Test_mdbm_OrdinaryReaplceData_StoreWithTryLockShared(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()

	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	for i := 0; i <= loopLimit; i++ {
		rv, err := dbm.StoreWithTryLockShared(i, i, mdbm.Replace)
		assert.AssertNil(t, err, "failured, return Value mismatch. value=%d, err=%v\n", rv, err)
	}
}

func Test_mdbm_OrdinaryReaplceData_StoreWithTryPlock(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()

	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	for i := 0; i <= loopLimit; i++ {
		rv, err := dbm.StoreWithTryPlock(i, i, mdbm.Replace, mdbm.Rdrw)
		assert.AssertNil(t, err, "failured, return Value mismatch. value=%d, err=%v\n", rv, err)
	}
}

func Test_mdbm_OrdinaryFetchData_Fetch(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()

	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	for i := 0; i <= loopLimit; i++ {
		val, err := dbm.Fetch(i)
		assert.AssertNil(t, err, "failured, return Value mismatch. value=%s, err=%v\n", val, err)
		assert.AssertEquals(t, strconv.Itoa(i), val, "return Value mismatch.\nExpected: %v\nActual: %v", i, val)
	}
}

func Test_mdbm_OrdinaryFetchData_RandomFetch(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()

	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	r1 := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i <= loopLimit; i++ {

		val, err := dbm.Fetch(r1.Intn(loopLimit))
		assert.AssertNil(t, err, "failured, return Value mismatch. value=%s, err=%v\n", val, err)
		assert.AssertEquals(t, strconv.Itoa(i), val, "return Value mismatch.\nExpected: %v\nActual: %v", i, val)
	}
}

func Test_mdbm_OrdinaryFetchData_FetchWithLock(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()

	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	for i := 0; i <= loopLimit; i++ {
		val, err := dbm.FetchWithLock(i)
		assert.AssertNil(t, err, "failured, return Value mismatch. value=%s, err=%v\n", val, err)
		assert.AssertEquals(t, strconv.Itoa(i), val, "return Value mismatch.\nExpected: %v\nActual: %v", i, val)
	}
}

func Test_mdbm_OrdinaryFetchData_FetchWithLockSmart(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()

	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	for i := 0; i <= loopLimit; i++ {
		val, err := dbm.FetchWithLockSmart(i, mdbm.Rdrw)
		assert.AssertNil(t, err, "failured, return Value mismatch. value=%s, err=%v\n", val, err)
		assert.AssertEquals(t, strconv.Itoa(i), val, "return Value mismatch.\nExpected: %v\nActual: %v", i, val)
	}
}

func Test_mdbm_OrdinaryFetchData_FetchWithLockShared(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()

	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	for i := 0; i <= loopLimit; i++ {
		val, err := dbm.FetchWithLockShared(i)
		assert.AssertNil(t, err, "failured, return Value mismatch. value=%s, err=%v\n", val, err)
		assert.AssertEquals(t, strconv.Itoa(i), val, "return Value mismatch.\nExpected: %v\nActual: %v", i, val)
	}
}

func Test_mdbm_OrdinaryFetchData_FetchWithPlock(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()

	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	for i := 0; i <= loopLimit; i++ {
		val, err := dbm.FetchWithPlock(i, mdbm.Rdrw)
		assert.AssertNil(t, err, "failured, return Value mismatch. value=%s, err=%v\n", val, err)
		assert.AssertEquals(t, strconv.Itoa(i), val, "return Value mismatch.\nExpected: %v\nActual: %v", i, val)
	}
}

func Test_mdbm_OrdinaryFetchData_Random_NonePreLoad_Fetch(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()

	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i <= loopLimit; i++ {

		val, err := dbm.Fetch(r.Intn(loopLimit))
		assert.AssertNil(t, err, "failured, return Value mismatch. value=%s, err=%v\n", val, err)
		assert.AssertEquals(t, strconv.Itoa(i), val, "return Value mismatch.\nExpected: %v\nActual: %v", i, val)
	}
}

func Test_mdbm_OrdinaryFetchData_Random_PreLoad_Fetch(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()
	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	rv, err := dbm.PreLoad()
	assert.AssertNil(t, err, "failured, can't pre-load the mdbm, path=%s, rv=%d, err=%v", pathTestDBM1, rv, err)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i <= loopLimit; i++ {

		val, err := dbm.Fetch(r.Intn(loopLimit))
		assert.AssertNil(t, err, "failured, return Value mismatch. value=%s, err=%v\n", val, err)
		assert.AssertEquals(t, strconv.Itoa(i), val, "return Value mismatch.\nExpected: %v\nActual: %v", i, val)
	}
}

func Test_mdbm_OrdinaryFetchData_Random_NonePreLoad_FetchWithLock(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()

	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i <= loopLimit; i++ {

		val, err := dbm.FetchWithLock(r.Intn(loopLimit))
		assert.AssertNil(t, err, "failured, return Value mismatch. value=%d, err=%v\n", val, err)
		assert.AssertEquals(t, strconv.Itoa(i), val, "return Value mismatch.\nExpected: %v\nActual: %v", i, val)
	}
}

func Test_mdbm_OrdinaryFetchData_Random_PreLoad_FetchWithLock(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()
	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	rv, err := dbm.PreLoad()
	assert.AssertNil(t, err, "failured, can't pre-load the mdbm, path=%s, rv=%d, err=%v", pathTestDBM1, rv, err)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i <= loopLimit; i++ {

		val, err := dbm.FetchWithLock(r.Intn(loopLimit))
		assert.AssertNil(t, err, "failured, return Value mismatch. value=%s, err=%v\n", val, err)
		assert.AssertEquals(t, strconv.Itoa(i), val, "return Value mismatch.\nExpected: %v\nActual: %v", i, val)
	}
}

func Test_mdbm_OrdinaryFetchData_Random_NonePreLoad_FetchWithLockSmart(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()

	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i <= loopLimit; i++ {

		val, err := dbm.FetchWithLockSmart(r.Intn(loopLimit), mdbm.Rdrw)
		assert.AssertNil(t, err, "failured, return Value mismatch. value=%s, err=%v\n", val, err)
		assert.AssertEquals(t, strconv.Itoa(i), val, "return Value mismatch.\nExpected: %v\nActual: %v", i, val)
	}
}

func Test_mdbm_OrdinaryFetchData_Random_PreLoad_FetchWithLockSmart(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()
	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	rv, err := dbm.PreLoad()
	assert.AssertNil(t, err, "failured, can't pre-load the mdbm, path=%s, rv=%d, err=%v", pathTestDBM1, rv, err)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i <= loopLimit; i++ {

		val, err := dbm.FetchWithLockSmart(r.Intn(loopLimit), mdbm.Rdrw)
		assert.AssertNil(t, err, "failured, return Value mismatch. value=%s, err=%v\n", val, err)
		assert.AssertEquals(t, strconv.Itoa(i), val, "return Value mismatch.\nExpected: %v\nActual: %v", i, val)
	}
}

func Test_mdbm_LockShared(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()
	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	rv, err := dbm.LockShared()
	assert.AssertEquals(t, 1, rv, "failured, Locks the database for shared access by readers, excluding access to writers., path=%s, err=%v", pathTestDBM1, err)
	assert.AssertNil(t, err, "failured, Locks the database for shared access by readers, excluding access to writers., path=%s, err=%v", pathTestDBM1, err)
}

func Test_mdbm_TryLockShared(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()
	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	rv, err := dbm.TryLockShared()
	assert.AssertEquals(t, 1, rv, "failured, locks the database for shared access by readers, excluding access to writers, path=%s, err=%v", pathTestDBM1, err)
	assert.AssertNil(t, err, "failured, locks the database for shared access by readers, excluding access to writers, path=%s, err=%v", pathTestDBM1, err)
}

func Test_mdbm_GetLockMode(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()
	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	rv, err := dbm.GetLockMode()
	assert.AssertEquals(t, 1, rv, "failured, gets the mdbm's lock mode, path=%s, err=%v", pathTestDBM1, err)
	assert.AssertNil(t, err, "failured, gets the mdbm's lock mode, path=%s, err=%v", pathTestDBM1, err)

}

func Test_mdbm_MutipleDataType_Store(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.Open(pathTestDBMLarge, mdbm.Create|mdbm.Rdrw|mdbm.LargeObjects, 0644, 0, 0)
	defer dbm.EasyClose()
	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	rv, err := dbm.PreLoad()
	assert.AssertNil(t, err, "failured, can't pre-load the mdbm, path=%s, rv=%d, err=%v", pathTestDBM1, rv, err)

	rv, err = dbm.StoreWithLock(true, time.Now().UnixNano(), mdbm.Replace)
	assert.AssertNil(t, err, "failured, return Value mismatch. value=%v, err=%v\n", rv, err)

	rv, err = dbm.StoreWithLock(false, time.Now().UnixNano(), mdbm.Replace)
	assert.AssertNil(t, err, "failured, return Value mismatch. value=%v, err=%v\n", rv, err)

	rv, err = dbm.StoreWithLock("true", time.Now().UnixNano(), mdbm.Replace)
	assert.AssertNil(t, err, "failured, return Value mismatch. value=%v, err=%v\n", rv, err)

	rv, err = dbm.StoreWithLock("false", time.Now().UnixNano(), mdbm.Replace)
	assert.AssertNil(t, err, "failured, return Value mismatch. value=%v, err=%v\n", rv, err)

	rv, err = dbm.StoreWithLock(byte(77), time.Now().UnixNano(), mdbm.Replace)
	assert.AssertNil(t, err, "failured, return Value mismatch. value=%v, err=%v\n", rv, err)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i <= loopLimit; i++ {

		rv, err = dbm.StoreWithLock(int8(r.Intn(100)), time.Now().UnixNano(), mdbm.Replace)
		assert.AssertNil(t, err, "failured, return Value mismatch. value=%v, err=%v\n", rv, err)

		rv, err = dbm.StoreWithLock(int16(r.Intn(100)), time.Now().UnixNano(), mdbm.Replace)
		assert.AssertNil(t, err, "failured, return Value mismatch. value=%v, err=%v\n", rv, err)

		rv, err = dbm.StoreWithLock(uint16(i), time.Now().UnixNano(), mdbm.Replace)
		assert.AssertNil(t, err, "failured, return Value mismatch. value=%v, err=%v\n", rv, err)

		rv, err = dbm.StoreWithLock(uint32(i), time.Now().UnixNano(), mdbm.Replace)
		assert.AssertNil(t, err, "failured, return Value mismatch. value=%v, err=%v\n", rv, err)

		rv, err = dbm.StoreWithLock(r.Int31(), time.Now().UnixNano(), mdbm.Replace)
		assert.AssertNil(t, err, "failured, return Value mismatch. value=%v, err=%v\n", rv, err)

		rv, err = dbm.StoreWithLock(r.Int63(), time.Now().UnixNano(), mdbm.Replace)
		assert.AssertNil(t, err, "failured, return Value mismatch. value=%v, err=%v\n", rv, err)

		rv, err = dbm.StoreWithLock(r.Uint32(), time.Now().UnixNano(), mdbm.Replace)
		assert.AssertNil(t, err, "failured, return Value mismatch. value=%v, err=%v\n", rv, err)

		rv, err = dbm.StoreWithLock(r.Float32(), time.Now().UnixNano(), mdbm.Replace)
		assert.AssertNil(t, err, "failured, return Value mismatch. value=%v, err=%v\n", rv, err)

		rv, err = dbm.StoreWithLock(r.Float64(), time.Now().UnixNano(), mdbm.Replace)
		assert.AssertNil(t, err, "failured, return Value mismatch. value=%v, err=%v\n", rv, err)

		rv, err = dbm.StoreWithLock(int64(r.Int()), time.Now().UnixNano(), mdbm.Replace)
		assert.AssertNil(t, err, "failured, return Value mismatch. value=%v, err=%v\n", rv, err)

		rv, err = dbm.StoreWithLock(uint64(r.Int()), time.Now().UnixNano(), mdbm.Replace)
		assert.AssertNil(t, err, "failured, return Value mismatch. value=%v, err=%v\n", rv, err)
	}
}

func Test_mdbm_DupHandle_AfterClose(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.Open(pathTestDBMLarge, mdbm.Create|mdbm.Rdrw|mdbm.LargeObjects, 0644, 0, 0)
	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	_, err = dbm.DupHandle()
	assert.AssertNil(t, err, "failured, a pointer of the Duplicate an existing database handle, path=%s, err=%v", pathTestDBM1, err)

	dbm.EasyClose()

	_, err = dbm.DupHandle()
	assert.AssertNotNil(t, err, "failured, return of closed db handler, err=%v", err)

}

func Test_mdbm_LogMinLevel_WrongOption(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()
	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	err = dbm.LogMinLevel(mdbm.LargeObjects)
	assert.AssertNotNil(t, err, "oops!. mdbm.LogMinLevel can't check a argument=%d", int(mdbm.LargeObjects))
}

func Test_mdbm_LogPlugin(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()
	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	err = dbm.LogPlugin(mdbm.LogToFile)
	assert.AssertNil(t, err, "failured, can't set logging to file, err=%v", err)

	err = dbm.LogPlugin(mdbm.LogToSkip)
	assert.AssertNil(t, err, "failured, can't set logging to /dev/null , err=%v", err)

	err = dbm.LogPlugin(mdbm.LogToStdErr)
	assert.AssertNil(t, err, "failured, can't set loging to /dev/stdout, err=%v", err)

	err = dbm.LogPlugin(mdbm.LogToSysLog)
	assert.AssertNil(t, err, "failured, can't set loging to syslog, err=%v", err)

	err = dbm.LogPlugin(mdbm.LargeObjects)
	assert.AssertNotNil(t, err, "oops!. mdbm.LogPlugin can't check a argument=%d", int(mdbm.LargeObjects))

}

func Test_mdbm_LogToAutoFile(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()
	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	err = dbm.LogPlugin(mdbm.LogToFile)
	assert.AssertNil(t, err, "failured, can't set logging to file, err=%v", err)

	rv, err := dbm.LogToAutoFile()
	assert.AssertNil(t, err, "failured, can't set logging to file, rv=%d, err=%v", rv, err)
}

func Test_mdbm_AnyDataType_Store(t *testing.T) {

	var rv int
	var err error

	dbm := mdbm.NewMDBM()
	err = dbm.EasyOpen(pathTestDBMAnyDataType, 0644)
	defer dbm.EasyClose()
	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBMAnyDataType, err)

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

	rv, err = dbm.StoreWithLock(vbyte, vbyte, mdbm.Insert)
	assert.AssertNil(t, err, "failed, can't data(byte) add to the mdbm file(=%s)\nrv=%d, err=%v", pathTestDBMAnyDataType, rv, err)

	rv, err = dbm.StoreWithLock(varbyte, varbyte, mdbm.Insert)
	assert.AssertNil(t, err, "failed, can't data([]byte) add to the mdbm file(=%s)\nrv=%d, err=%v", pathTestDBMAnyDataType, rv, err)

	rv, err = dbm.StoreWithLock(vstring, vstring, mdbm.Insert)
	assert.AssertNil(t, err, "failed, can't data(string) add to the mdbm file(=%s)\nrv=%d, err=%v", pathTestDBMAnyDataType, rv, err)

	rv, err = dbm.StoreWithLock(vint, vint, mdbm.Insert)
	assert.AssertNil(t, err, "failed, can't data(int) add to the mdbm file(=%s)\nrv=%d, err=%v", pathTestDBMAnyDataType, rv, err)

	rv, err = dbm.StoreWithLock(vint16, vint16, mdbm.Insert)
	assert.AssertNil(t, err, "failed, can't data(int16) add to the mdbm file(=%s)\nrv=%d, err=%v", pathTestDBMAnyDataType, rv, err)

	rv, err = dbm.StoreWithLock(vint32, vint32, mdbm.Insert)
	assert.AssertNil(t, err, "failed, can't data(int32) add to the mdbm file(=%s)\nrv=%d, err=%v", pathTestDBMAnyDataType, rv, err)

	rv, err = dbm.StoreWithLock(vint64, vint64, mdbm.Insert)
	assert.AssertNil(t, err, "failed, can't data(int32) add to the mdbm file(=%s)\nrv=%d, err=%v", pathTestDBMAnyDataType, rv, err)

	rv, err = dbm.StoreWithLock(vuint, vuint, mdbm.Insert)
	assert.AssertNil(t, err, "failed, can't data(uint) add to the mdbm file(=%s)\nrv=%d, err=%v", pathTestDBMAnyDataType, rv, err)

	rv, err = dbm.StoreWithLock(vuint16, vuint16, mdbm.Insert)
	assert.AssertNil(t, err, "failed, can't data(uint16) add to the mdbm file(=%s)\nrv=%d, err=%v", pathTestDBMAnyDataType, rv, err)

	rv, err = dbm.StoreWithLock(vuint32, vuint32, mdbm.Insert)
	assert.AssertNil(t, err, "failed, can't data(uint32) add to the mdbm file(=%s)\nrv=%d, err=%v", pathTestDBMAnyDataType, rv, err)

	rv, err = dbm.StoreWithLock(vuint64, vuint64, mdbm.Insert)
	assert.AssertNil(t, err, "failed, can't data(uint32) add to the mdbm file(=%s)\nrv=%d, err=%v", pathTestDBMAnyDataType, rv, err)

	rv, err = dbm.StoreWithLock(vfloat32, vfloat32, mdbm.Insert)
	assert.AssertNil(t, err, "failed, can't data(float32) add to the mdbm file(=%s)\nrv=%d, err=%v", pathTestDBMAnyDataType, rv, err)

	rv, err = dbm.StoreWithLock(vfloat64, vfloat64, mdbm.Insert)
	assert.AssertNil(t, err, "failed, can't data(float64) add to the mdbm file(=%s)\nrv=%d, err=%v", pathTestDBMAnyDataType, rv, err)

	rv, err = dbm.StoreWithLock(vcomplex64, vcomplex64, mdbm.Insert)
	assert.AssertNil(t, err, "failed, can't data(complex64) add to the mdbm file(=%s)\nrv=%d, err=%v", pathTestDBMAnyDataType, rv, err)

	rv, err = dbm.StoreWithLock(vcomplex128, vcomplex128, mdbm.Insert)
	assert.AssertNil(t, err, "failed, can't data(complex128) add to the mdbm file(=%s)\nrv=%d, err=%v", pathTestDBMAnyDataType, rv, err)

	rv, err = dbm.Sync()
	assert.AssertNil(t, err, "failured, can't sync database, path=%s, rv=%d, err=%v", pathTestDBMAnyDataType, rv, err)

	//validation
	var cnt int
	key, val, err := dbm.First()
	assert.AssertNil(t, err, "failured, can't get a first records, path=%s, err=%v", pathTestDBMAnyDataType, err)
	assert.AssertEquals(t, key, val, "key and value mismatch.\nKey=%s, Value=%s", key, val)

	for {

		key, val, err := dbm.Next()
		if err != nil {
			log.Fatalf("key=%s, val=%s, err=%v", key, val, err)
		}

		if len(key) < 1 {
			break
		}

		assert.AssertEquals(t, key, val, "key and value mismatch.\nKey=%s, Value=%s", key, val)
		cnt++
	}

	assert.AssertEquals(t, cnt, 15, "count of records value mismatch.\ngot=%d, want=%d", cnt, 15)
}

func Test_mdbm_TryLock_UnLock(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	defer dbm.EasyClose()
	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBM1, err)

	rv, err := dbm.TryLock()
	assert.AssertNil(t, err, "failured, can't try-locking, path=%s, rv=%d, err=%v", pathTestDBM1, rv, err)

	rv, err = dbm.Unlock()
	assert.AssertNil(t, err, "failured, can't un-locking, path=%s, rv=%d, err=%v", pathTestDBM1, rv, err)
}

func Test_mdbm_FetchInfo(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBMLarge, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	iter := dbm.GetNewIter()

	uint32zero := uint32(0)
	intzero := int(0)

	for i := 0; i <= 1; i++ {

		var retval string

		rv, copiedval, info, goiter, err := dbm.FetchInfo(i, &retval, &iter)
		assert.AssertNil(t, err, "failured, can't get fetch infor, path=%s, rv=%d, err=%v", pathTestDBM1, rv, err)
		assert.AssertEquals(t, string(i), copiedval, "failured, return Value mismatch.\nExpected: %v\nActual: %v", i, copiedval)

		assert.AssertEquals(t, info.Flags, uint32zero, "failured, return Value mismatch.\nExpected: %v\nActual: %v", info.Flags, uint32zero)
		assert.AssertEquals(t, info.CacheNumAccesses, uint32zero, "failured, return Value mismatch.\nExpected: %v\nActual: %v", info.CacheNumAccesses, uint32zero)
		assert.AssertEquals(t, info.CacheAccessTime, uint32zero, "failured, return Value mismatch.\nExpected: %v\nActual: %v", info.CacheAccessTime, uint32zero)

		assert.AssertNotEquals(t, goiter.PageNo, uint32zero, "failured, return Value mismatch.\nExpected: %d\nActual: %d", goiter.PageNo, uint32zero)
		assert.AssertNotEquals(t, goiter.Next, intzero, "failured, return Value mismatch.\nExpected: %d\nActual: %d", goiter.Next, intzero)
	}
}

func Test_mdbm_DeleteWithLock(t *testing.T) {

	var rv int

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBMDelete, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	for i := 0; i <= loopLimit; i++ {
		rv, err = dbm.StoreWithLock(i, time.Now().UnixNano(), mdbm.Replace)
		assert.AssertNil(t, err, "failured, return Value mismatch. value=%v, err=%v\n", rv, err)
	}

	rv, err = dbm.Sync()
	assert.AssertNil(t, err, "failured, mdbm.Sync(). rv=%v, err=%v\n", rv, err)

	for i := 0; i <= loopLimit; i++ {
		rv, err = dbm.DeleteWithLock(i)
		assert.AssertNil(t, err, "failured, can't delete record, return value=%v, err=%v\n", rv, err)
	}

	rv, err = dbm.Sync()
	assert.AssertNil(t, err, "failured, mdbm.Sync(). rv=%v, err=%v\n", rv, err)

	for i := 0; i <= loopLimit; i++ {

		val, err := dbm.Fetch(i)
		assert.AssertNotNil(t, err, "failured, can't delete record, value=%v, err=%v\n", val, err)
	}
}

func Test_mdbm_EasyGetNumOfRows(t *testing.T) {
	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBMLarge, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	cnt, err := dbm.EasyGetNumOfRows()

	assert.AssertNil(t, err, "failured, can't obtain the count of number of rows, err=%v\n", err)
	assert.AssertGreaterThanEqualTo(t, cnt, uint64(1), "failured, return Value mismatch.\nExpected: >=%d\nActual: %d", 1, cnt)

	dbm.EasyClose()
	_, err = dbm.EasyGetNumOfRows()
	assert.AssertNotNil(t, err, "failured, can't check the mdbm closed, err=%v\n", err)

}

func Test_mdbm_EasyGetKeyList(t *testing.T) {
	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBMLarge, 0644)
	if err != nil {
		log.Fatalf("failed mdbm.EasyOpen(), err=%v", err)
	}
	defer dbm.EasyClose()

	keys, err := dbm.EasyGetKeyList()

	assert.AssertNil(t, err, "failured, can't obtain the list of key, err=%v\n", err)
	assert.AssertGreaterThanEqualTo(t, len(keys), uint64(1), "failured, return Value mismatch.\nExpected: >=%d\nActual: %d", 1, len(keys))

	dbm.EasyClose()
	_, err = dbm.EasyGetKeyList()
	assert.AssertNotNil(t, err, "failured, can't check the mdbm closed, err=%v\n", err)
}

func Test_mdbm_LockSmart_UnLockSmart(t *testing.T) {

	var rv int
	var err error

	dbm := mdbm.NewMDBM()
	err = dbm.EasyOpen(pathTestDBMAnyDataType, 0644)
	defer dbm.EasyClose()
	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBMAnyDataType, err)

	vint := int(123456789)

	rv, err = dbm.LockSmart(vint, mdbm.Rdrw)
	assert.AssertNil(t, err, "failured, can't obtain TryLockSmart, rv=%d ,err=%v\n", rv, err)

	rv, err = dbm.Store(vint, vint, mdbm.Replace)
	assert.AssertNil(t, err, "failed, can't data(=%d) add to the mdbm file(=%s)\nrv=%d, err=%v", vint, pathTestDBMAnyDataType, rv, err)

	rv, err = dbm.UnLockSmart(vint, mdbm.Rdrw)
	assert.AssertNil(t, err, "failured, can't obtain TryLockSmart, rv=%d ,err=%v\n", rv, err)
}

func Test_mdbm_TryLockSmart_UnLockSmart(t *testing.T) {

	var rv int
	var err error

	dbm := mdbm.NewMDBM()
	err = dbm.EasyOpen(pathTestDBMAnyDataType, 0644)
	defer dbm.EasyClose()
	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBMAnyDataType, err)

	vint := int(123456789)

	rv, err = dbm.TryLockSmart(vint, mdbm.Rdrw)
	assert.AssertNil(t, err, "failured, can't obtain TryLockSmart, rv=%d ,err=%v\n", rv, err)

	rv, err = dbm.Store(vint, vint, mdbm.Replace)
	assert.AssertNil(t, err, "failed, can't data(=%d) add to the mdbm file(=%s)\nrv=%d, err=%v", vint, pathTestDBMAnyDataType, rv, err)

	rv, err = dbm.UnLockSmart(vint, mdbm.Rdrw)
	assert.AssertNil(t, err, "failured, can't obtain TryLockSmart, rv=%d ,err=%v\n", rv, err)
}

func Test_mdbm_Truncate(t *testing.T) {

	var rv int
	var val string
	var err error

	dbm := mdbm.NewMDBM()
	err = dbm.EasyOpen(pathTestDBMDelete, 0644)
	defer dbm.EasyClose()
	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBMAnyDataType, err)

	for i := 0; i <= loopLimit; i++ {
		rv, err = dbm.StoreWithLock(i, time.Now().UnixNano(), mdbm.Replace)
		assert.AssertNil(t, err, "failured, return Value mismatch. value=%v, err=%v\n", rv, err)
	}

	rv, err = dbm.Sync()
	assert.AssertNil(t, err, "failured, mdbm.Sync(). rv=%v, err=%v\n", rv, err)

	err = dbm.Truncate()
	assert.AssertNil(t, err, "failured, can't truncate mdbm, err=%v\n", err)

	rv, err = dbm.Sync()
	assert.AssertNil(t, err, "failured, mdbm.Sync(). rv=%v, err=%v\n", rv, err)

	for i := 0; i <= loopLimit; i++ {

		val, err = dbm.Fetch(i)
		assert.AssertNotNil(t, err, "failured, can't delete record, value=%v, err=%v\n", val, err)
	}
}

func Test_mdbm_GetHashValue(t *testing.T) {

	var rv uint32
	var err error

	dbm := mdbm.NewMDBM()
	defer dbm.EasyClose()

	rv, _ = dbm.GetHashValue(1, mdbm.HashCRC32)
	assert.AssertEquals(t, rv, uint32(2667302803), "return Value mismatch.\nExpected: %v\nActual: %v", rv, uint32(2667302803))

	rv, _ = dbm.GetHashValue(1, mdbm.HashEJB)
	assert.AssertEquals(t, rv, uint32(17), "return Value mismatch.\nExpected: %v\nActual: %v", rv, uint32(17))

	rv, _ = dbm.GetHashValue(1, mdbm.HashPHONG)
	assert.AssertEquals(t, rv, uint32(2621031278), "return Value mismatch.\nExpected: %v\nActual: %v", rv, uint32(2621031278))

	rv, _ = dbm.GetHashValue(1, mdbm.HashOZ)
	assert.AssertEquals(t, rv, uint32(49), "return Value mismatch.\nExpected: %v\nActual: %v", rv, uint32(49))

	rv, _ = dbm.GetHashValue(1, mdbm.HashTOREK)
	assert.AssertEquals(t, rv, uint32(49), "return Value mismatch.\nExpected: %v\nActual: %v", rv, uint32(49))

	rv, _ = dbm.GetHashValue(1, mdbm.HashFNV)
	assert.AssertEquals(t, rv, uint32(1224750888), "return Value mismatch.\nExpected: %v\nActual: %v", rv, uint32(1224750888))

	rv, _ = dbm.GetHashValue(1, mdbm.HashSTL)
	assert.AssertEquals(t, rv, uint32(49), "return Value mismatch.\nExpected: %v\nActual: %v", rv, uint32(49))

	rv, _ = dbm.GetHashValue(1, mdbm.HashMD5)
	assert.AssertEquals(t, rv, uint32(943901380), "return Value mismatch.\nExpected: %v\nActual: %v", rv, uint32(943901380))

	rv, _ = dbm.GetHashValue(1, mdbm.HashSHA1)
	assert.AssertEquals(t, rv, uint32(723085877), "return Value mismatch.\nExpected: %v\nActual: %v", rv, uint32(723085877))

	rv, _ = dbm.GetHashValue(1, mdbm.HashJENKINS)
	assert.AssertEquals(t, rv, uint32(2366665294), "return Value mismatch.\nExpected: %v\nActual: %v", rv, uint32(2366665294))

	rv, _ = dbm.GetHashValue(1, mdbm.HashHSIEH)
	assert.AssertEquals(t, rv, uint32(3927678806), "return Value mismatch.\nExpected: %v\nActual: %v", rv, uint32(3927678806))

	rv, _ = dbm.GetHashValue(1, mdbm.MaxHash)
	assert.AssertEquals(t, rv, uint32(3927678806), "return Value mismatch.\nExpected: %v\nActual: %v", rv, uint32(3927678806))

	rv, _ = dbm.GetHashValue(1, mdbm.DefaultHash)
	assert.AssertEquals(t, rv, uint32(1224750888), "return Value mismatch.\nExpected: %v\nActual: %v", rv, uint32(1224750888))

	rv, err = dbm.GetHashValue(1, mdbm.DefaultHash)
	assert.AssertNotEquals(t, rv, uint32(1), "return Value mismatch.\nrv=%d, err=%v", rv, err)

	_, err = dbm.GetHashValue(1, mdbm.LargeObjects)
	assert.AssertNotNil(t, err, "failured, can't check hash type err=%v", err)
}

func Test_mdbm_StoreStr(t *testing.T) {

	var rv int
	var err error

	dbm := mdbm.NewMDBM()
	err = dbm.EasyOpen(pathTestDBMStr, 0644)
	defer dbm.EasyClose()
	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBMStr, err)

	rv, err = dbm.SetCacheMode(mdbm.CacheModeMax)
	assert.AssertNil(t, err, "failured, can't set cachemode(=%s), path=%s, rv=%d, err=%v", pathTestDBMStr, rv, err)

	for i := 0; i <= 10; i++ {
		//for i := 0; i <= loopLimit; i++ {
		rv, err = dbm.StoreStr(i, i, mdbm.Insert)
		assert.AssertNil(t, err, "failed, can't data(=%d) add to the mdbm file(=%s)\nrv=%d, err=%v", i, pathTestDBMStr, rv, err)
	}

	key, val, err := dbm.First()
	assert.AssertNil(t, err, "failured, can't get first recods, path=%s, err=%v", pathTestDBMStr, err)

	log.Println(key, val)
	for {

		key, val, err = dbm.Next()
		if err != nil {
			log.Fatalf("key=%s, val=%s, err=%v", key, val, err)
		}

		if len(key) < 1 {
			break
		}

		log.Println(key, val)
	}

	rv, err = dbm.Sync()
	assert.AssertNil(t, err, "failured, mdbm.Sync(). rv=%v, err=%v\n", rv, err)
}

func Test_mdbm_Double_Close(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBM1, 0644)
	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBMStr, err)
	defer dbm.EasyClose()
	dbm.EasyClose()
}

func Test_mdbm_SetHash(t *testing.T) {

	var rv int
	var err error

	dbm := mdbm.NewMDBM()
	err = dbm.EasyOpen(pathTestDBMHash, 0644)
	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBMHash, err)
	defer dbm.EasyClose()

	hashlist := []int{
		mdbm.HashCRC32,
		mdbm.HashEJB,
		mdbm.HashPHONG,
		mdbm.HashOZ,
		mdbm.HashTOREK,
		mdbm.HashFNV,
		mdbm.HashSTL,
		mdbm.HashMD5,
		mdbm.HashSHA1,
		mdbm.HashJENKINS,
		mdbm.HashHSIEH,
		mdbm.MaxHash,
		mdbm.DefaultHash,
	}

	for _, hashtype := range hashlist {

		err = dbm.SetHash(mdbm.HashOZ)
		assert.AssertNil(t, err, "failured, can't set hash(=%d) to the mdbm, path=%s, err=%v", hashtype, pathTestDBMHash, err)

		rv, err = dbm.GetHash()
		assert.AssertNil(t, err, "failured, can't get hash type of the mdbm, path=%s, err=%v", pathTestDBMHash, err)
		assert.AssertEquals(t, rv, hashtype, "return Value mismatch.\nExpected: %v\nActual: %v", hashtype, rv)
	}

	err = dbm.SetHash(mdbm.LargeObjects)
	assert.AssertNil(t, err, "failured, can't check wrong option, path=%s", pathTestDBMHash)
}

func Test_mdbm_ReplaceDB(t *testing.T) {

	dbm := mdbm.NewMDBM()
	err := dbm.EasyOpen(pathTestDBMLarge, 0644)
	assert.AssertNil(t, err, "failured, can't open the mdbm, path=%s, err=%v", pathTestDBMLarge, err)
	defer dbm.EasyClose()

	err = dbm.ReplaceDB(pathTestDBMReplace)
	assert.AssertNil(t, err, "failured, can't replace %s to %s, err=%v", pathTestDBMLarge, pathTestDBMReplace, err)
}
